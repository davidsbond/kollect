package agent_test

import (
	"context"

	"github.com/davidsbond/kollect/internal/event"
)

type (
	MockEventWriter struct {
		emitted   chan bool
		clusterID string

		event event.Event
	}
)

func (m *MockEventWriter) Write(_ context.Context, evt event.Event) error {
	m.event = evt
	m.emitted <- true
	return nil
}

func (m *MockEventWriter) Wait() {
	<-m.emitted
	return
}
