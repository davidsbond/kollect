// Package metrics provides functions for serving and registering prometheus metrics.
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Serve prometheus metrics via HTTP on the provided router.
func Serve(r *http.ServeMux) {
	r.Handle("/__/metrics", promhttp.Handler())
}

// Register a prometheus collector.
func Register(collectors ...prometheus.Collector) {
	prometheus.DefaultRegisterer.MustRegister(collectors...)
}
