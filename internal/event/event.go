// Package event contains utilities for interacting with various event-stream providers. Including the ability
// to write and read from event-streaming sources.
package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/kollect/proto/kollect/event/v1"

	// Gocloud driver for GCP PubSub.
	_ "gocloud.dev/pubsub/gcppubsub"
	// Gocloud driver for in-memory event streaming.
	_ "gocloud.dev/pubsub/mempubsub"
	// Gocloud driver for Apache Kafka.
	_ "gocloud.dev/pubsub/kafkapubsub"
	// Gocloud driver for NATS.
	_ "gocloud.dev/pubsub/natspubsub"
	// Gocloud driver for Azure Service Bus.
	_ "gocloud.dev/pubsub/azuresb"
	// Gocloud driver for RabbitMQ.
	_ "gocloud.dev/pubsub/rabbitpubsub"
	// Gocloud driver for AWS SNS/SQS.
	_ "gocloud.dev/pubsub/awssnssqs"
)

type (
	// The Event type describes something that has happened at a particular point in time.
	Event struct {
		ID        string
		Key       string
		Timestamp time.Time
		AppliesAt time.Time
		Payload   proto.Message
	}

	// The Handler type is a function that processes an inbound event.
	Handler func(ctx context.Context, evt Event) error

	// The Option type is a function that can modify an event value.
	Option func(e *Event)
)

// New returns a new Event instance that contains the provided payload. A unique identifier
// is generated and timestamps are set to now.
func New(payload proto.Message, opts ...Option) Event {
	e := Event{
		ID:        uuid.Must(uuid.NewUUID()).String(),
		Timestamp: time.Now(),
		AppliesAt: time.Now(),
		Payload:   payload,
	}

	for _, opt := range opts {
		opt(&e)
	}

	return e
}

func (e Event) typeName() string {
	return string(e.Payload.ProtoReflect().Descriptor().FullName())
}

func (e Event) marshal() ([]byte, error) {
	any, err := anypb.New(e.Payload)
	if err != nil {
		return nil, err
	}

	envelope := &event.Envelope{
		Id:        e.ID,
		Timestamp: timestamppb.New(e.Timestamp),
		AppliesAt: timestamppb.New(e.AppliesAt),
		Payload:   any,
	}

	return proto.Marshal(envelope)
}

func unmarshal(b []byte) (Event, error) {
	var env event.Envelope
	if err := proto.Unmarshal(b, &env); err != nil {
		return Event{}, err
	}

	payload, err := env.Payload.UnmarshalNew()
	if err != nil {
		return Event{}, err
	}

	return Event{
		ID:        env.GetId(),
		Timestamp: env.GetTimestamp().AsTime(),
		AppliesAt: env.GetAppliesAt().AsTime(),
		Payload:   payload,
	}, nil
}

// WithAppliesAt returns an Option that can be provided to New to set the Event.AppliesAt field.
func WithAppliesAt(appliesAt time.Time) Option {
	return func(e *Event) {
		e.AppliesAt = appliesAt
	}
}

// WithKey returns an Option that can be provided to New to set the Event.Key field.
func WithKey(key string) Option {
	return func(e *Event) {
		e.Key = key
	}
}
