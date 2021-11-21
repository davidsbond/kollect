package event

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "kollect"
	subsystem = "events"
)

func init() {
	prometheus.MustRegister(eventsWritten, eventsRead, eventsIgnored)
}

var (
	eventsWritten = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "written_total",
		Help:      "Total number of events written to the stream",
	}, []string{"key", "type"})

	eventsRead = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "read_total",
		Help:      "Total number of events read from the stream",
	}, []string{"key", "type"})

	eventsIgnored = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "ignored_total",
		Help:      "Total number of events ignored from the stream",
	}, []string{"key", "type"})
)
