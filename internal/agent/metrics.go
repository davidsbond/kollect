package agent

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "kollect"
	subsystem = "resource"
)

func init() {
	prometheus.MustRegister(
		resourceCreated,
		resourceUpdated,
		resourceDeleted,
	)
}

var (
	resourceCreated = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "created_total",
		Help:      "Total number of resources created",
	}, []string{"group", "version", "kind", "namespace"})

	resourceUpdated = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "updated_total",
		Help:      "Total number of resources updated",
	}, []string{"group", "version", "kind", "namespace"})

	resourceDeleted = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "deleted_total",
		Help:      "Total number of resources deleted",
	}, []string{"group", "version", "kind", "namespace"})
)
