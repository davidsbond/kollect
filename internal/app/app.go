// Package app provides functions for creating command-line applications
package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/davidsbond/kollect/internal/closers"
	"github.com/davidsbond/kollect/internal/diagnostics"
	"github.com/davidsbond/kollect/internal/environment"
	"github.com/davidsbond/kollect/internal/flag"
	"github.com/davidsbond/kollect/internal/health"
	"github.com/davidsbond/kollect/internal/metrics"
	"github.com/davidsbond/kollect/internal/telemetry"
)

type (
	// The App type represents a command-line application.
	App struct {
		inner *cli.App
	}

	// The Option type is a function that can apply configuration to the command-line
	// application.
	Option func(app *cli.App)

	// The RunFunc type describes a method invoked to start a cli command.
	RunFunc func(ctx context.Context) error
)

// New returns a new App that functions as a command-line application configured via the provided Option functions.
func New(opts ...Option) *App {
	app := cli.NewApp()
	app.Version = environment.Version
	app.Name = environment.ApplicationName
	app.Usage = environment.ApplicationDescription
	app.Compiled = environment.Compiled()

	for _, opt := range opts {
		opt(app)
	}

	return &App{inner: app}
}

// Run the application.
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(environment.NewContext())
	defer cancel()

	action := a.inner.Action
	a.inner.Action = func(c *cli.Context) error {
		environment.SetMaxProcsToCPUQuota()

		svr := createOperationalServer()
		defer closers.Close(svr)

		tracer, err := telemetry.NewTracer(ctx)
		if err != nil {
			return err
		}
		defer closers.Close(tracer)

		return action(c)
	}

	return a.inner.RunContext(ctx, os.Args)
}

// WithRunner sets the action function to be used by the command-line application when Run is
// called.
func WithRunner(run RunFunc) Option {
	return func(app *cli.App) {
		app.Action = func(c *cli.Context) error {
			return run(c.Context)
		}
	}
}

// WithFlags sets command-line flags that can be applied before Run is called.
func WithFlags(flags ...flag.Flag) Option {
	return func(app *cli.App) {
		for _, fl := range flags {
			app.Flags = append(app.Flags, fl.Unwrap())
		}

		app.Flags = append(app.Flags, environment.Flags.Unwrap()...)
		app.Flags = append(app.Flags, telemetry.Flags.Unwrap()...)
		app.Flags = append(app.Flags, diagnostics.Flags.Unwrap()...)

		sort.Sort(cli.FlagsByName(app.Flags))
	}
}

func createOperationalServer() *http.Server {
	mux := http.NewServeMux()
	svr := &http.Server{Addr: ":8081", Handler: mux}

	metrics.Serve(mux)
	health.Serve(mux)
	diagnostics.Serve(mux)

	go func() {
		err := svr.ListenAndServe()
		switch {
		case errors.Is(err, http.ErrServerClosed):
			break
		case err != nil:
			log.Println(err)
		}
	}()

	return svr
}
