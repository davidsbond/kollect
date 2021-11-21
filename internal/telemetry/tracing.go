// Package telemetry provides functions for application tracing. Primarily used while debugging.
package telemetry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/davidsbond/kollect/internal/closers"
	"github.com/davidsbond/kollect/internal/environment"
)

// Span is a wrapper for trace.Span.
type Span struct {
	trace.Span
}

// SetAttributes is a wrapper around trace.Span.SetAttributes that uses an interface map to set attributes.
func (s Span) SetAttributes(m map[string]interface{}) {
	attrs := make([]attribute.KeyValue, len(m))
	i := 0
	for k, v := range m {
		attrs[i] = anyAttribute(k, v)
		i++
	}

	s.Span.SetAttributes(attrs...)
}

var (
	provider *sdk.TracerProvider

	// ErrUnsupportedProvider is the error given when a trace provider cannot be found for the provided
	// configuration.
	ErrUnsupportedProvider = errors.New("unsupported trace provider")

	// Discard is a sentinel error that is used to discard any current traces.
	// nolint: revive,stylecheck
	Discard = errors.New("discard")
)

// NewTracer constructs a tracer based on the telemetry flag and returns an io.Closer implementation used to stop
// the tracer. Tracers are configured via a URL where the scheme denotes the tracer to use. Currently noop, log, zipkin
// and jaeger are supported.
func NewTracer(ctx context.Context) (io.Closer, error) {
	u, err := url.Parse(telemetryURL)
	if err != nil {
		return nil, err
	}

	var exporter sdk.SpanExporter
	switch u.Scheme {
	case "noop":
		return closers.Noop, nil
	case "jaeger":
		exporter, err = jaeger.New(
			jaeger.WithAgentEndpoint(
				jaeger.WithAgentHost(u.Hostname()),
				jaeger.WithAgentPort(u.Port()),
			),
		)
	case "zipkin":
		u.Scheme = "http"
		exporter, err = zipkin.New(u.String())
	default:
		err = fmt.Errorf("%w: %s", ErrUnsupportedProvider, u.Scheme)
	}
	if err != nil {
		return nil, err
	}

	rs, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithFromEnv(),
		resource.WithDetectors(
			resource.StringDetector(semconv.SchemaURL, semconv.ServiceNameKey, func() (string, error) {
				return environment.ApplicationName, nil
			}),
			resource.StringDetector(semconv.SchemaURL, semconv.ServiceVersionKey, func() (string, error) {
				return environment.Version, nil
			}),
		),
	)
	if err != nil {
		return nil, err
	}

	provider = sdk.NewTracerProvider(
		sdk.WithResource(rs),
		sdk.WithSpanProcessor(sdk.NewSimpleSpanProcessor(exporter)),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(&propagation.TraceContext{})

	return closers.CloseFunc(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		return provider.Shutdown(ctx)
	}), nil
}

// WithinSpan wraps fn with a Span.
func WithinSpan(ctx context.Context, name string, fn func(ctx context.Context, span Span) error) error {
	ctx, span := otel.Tracer("").Start(ctx, name)

	err := fn(ctx, Span{span})
	switch {
	case errors.Is(err, Discard):
		return nil
	case err != nil:
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return err
	default:
		span.End()
		return nil
	}
}

type (
	mapCarrier map[string]string
)

func (m mapCarrier) Get(key string) string {
	return m[key]
}

func (m mapCarrier) Set(key string, value string) {
	m[key] = value
}

func (m mapCarrier) Keys() []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// Extract the tracer headers from the provided context as a map[string]string.
func Extract(ctx context.Context) map[string]string {
	m := mapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, m)
	return m
}

// Inject the tracer headers from the provided map into a new context.Context.
func Inject(ctx context.Context, m map[string]string) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, mapCarrier(m))
}

type (
	// The TracedMutex type is a wrapper around a sync.Mutex that creates a span when Lock is called.
	TracedMutex struct {
		mux *sync.Mutex
	}
)

// NewTracedMutex returns a new TracedMutex initialised with a sync.Mutex.
func NewTracedMutex() *TracedMutex {
	return &TracedMutex{mux: &sync.Mutex{}}
}

// Lock the mutex.
func (t *TracedMutex) Lock(ctx context.Context) {
	_, span := otel.Tracer("").Start(ctx, "Mutex.Lock")
	defer span.End()

	t.mux.Lock()
}

// Unlock the mutex.
func (t *TracedMutex) Unlock() {
	t.mux.Unlock()
}

func anyAttribute(k string, value interface{}) attribute.KeyValue {
	if value == nil {
		return attribute.String(k, "<nil>")
	}

	if stringer, ok := value.(fmt.Stringer); ok {
		return attribute.String(k, stringer.String())
	}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Array:
		rv = rv.Slice(0, rv.Len())
		fallthrough
	case reflect.Slice:
		switch reflect.TypeOf(value).Elem().Kind() {
		case reflect.Bool:
			return attribute.BoolSlice(k, rv.Interface().([]bool))
		case reflect.Int:
			return attribute.IntSlice(k, rv.Interface().([]int))
		case reflect.Int64:
			return attribute.Int64Slice(k, rv.Interface().([]int64))
		case reflect.Float64:
			return attribute.Float64Slice(k, rv.Interface().([]float64))
		case reflect.String:
			return attribute.StringSlice(k, rv.Interface().([]string))
		default:
			return attribute.String(k, "INVALID")
		}
	case reflect.Bool:
		return attribute.Bool(k, rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return attribute.Int64(k, rv.Int())
	case reflect.Float64:
		return attribute.Float64(k, rv.Float())
	case reflect.String:
		return attribute.String(k, rv.String())
	}
	if b, err := json.Marshal(value); b != nil && err == nil {
		return attribute.String(k, string(b))
	}
	return attribute.String(k, fmt.Sprint(value))
}
