// Package health provides functions for reporting the health status of individual components in the application.
package health

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/davidsbond/kollect/internal/environment"
)

type (
	// The Health type contains fields that describe the health of an application and its components.
	Health struct {
		Application string    `json:"application,omitempty"`
		Version     string    `json:"version,omitempty"`
		Compiled    time.Time `json:"compiled,omitempty"`
		Description string    `json:"description,omitempty"`
		Healthy     bool      `json:"healthy"`
		Checks      []Check   `json:"checks,omitempty"`
	}

	// The Check type describes the health of an individual component.
	Check struct {
		Name    string `json:"name"`
		Healthy bool   `json:"healthy"`
		Message string `json:"message,omitempty"`
	}

	// The ReadyFunc type is used to indicate if a component of the application is ready or not.
	ReadyFunc func() bool

	healthChecker struct {
		Name string
		Func func() error
	}
)

var (
	healthChecks = make([]healthChecker, 0)
	readyChecks  = make([]ReadyFunc, 0)
)

// AddHealthCheck adds a new named component to the health check results.
func AddHealthCheck(name string, check func() error) {
	healthChecks = append(healthChecks, healthChecker{
		Name: name,
		Func: check,
	})
}

// AddReadyCheck adds a new ReadyFunc to the readiness check. The readiness endpoint returns an error if one of
// the functions returns false.
func AddReadyCheck(fn ReadyFunc) {
	readyChecks = append(readyChecks, fn)
}

var errUnhealthy = errors.New("unhealthy")

// CheckKubernetesAPI returns a function that returns a non-nil error when the connection to the Kubernetes cluster
// is deemed unhealthy.
func CheckKubernetesAPI(cnf *rest.Config) func() error {
	return func() error {
		cl, err := kubernetes.NewForConfig(cnf)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		content, err := cl.Discovery().RESTClient().Get().AbsPath("/healthz").DoRaw(ctx)
		if err != nil {
			return err
		}

		if string(content) != "ok" {
			return fmt.Errorf("%w: %s", errUnhealthy, string(content))
		}

		return nil
	}
}

// Serve health/ready checks via HTTP on the provided router.
func Serve(r *http.ServeMux) {
	r.HandleFunc("/__/ready", func(w http.ResponseWriter, r *http.Request) {
		for _, fn := range readyChecks {
			if !fn() {
				http.Error(w, "not ready", http.StatusInternalServerError)
				return
			}
		}
	})

	r.HandleFunc("/__/health", func(w http.ResponseWriter, r *http.Request) {
		health := Health{
			Application: environment.ApplicationName,
			Version:     environment.Version,
			Compiled:    environment.Compiled(),
			Description: environment.ApplicationDescription,
			Checks:      make([]Check, len(healthChecks)),
			Healthy:     true,
		}

		for i, checker := range healthChecks {
			check := Check{
				Name:    checker.Name,
				Healthy: true,
			}

			if err := checker.Func(); err != nil {
				check.Healthy = false
				check.Message = err.Error()
				health.Healthy = false
			}

			health.Checks[i] = check
		}

		status := http.StatusOK
		if !health.Healthy {
			status = http.StatusInternalServerError
		}

		w.WriteHeader(status)
		enc := json.NewEncoder(w)
		enc.SetIndent("", strings.Repeat(" ", 4))
		if err := enc.Encode(health); err != nil {
			log.Println(err)
		}
	})
}
