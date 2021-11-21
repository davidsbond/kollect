// Package kollect provides types/functions for reading kollect resource events from supported event buses.
package kollect

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/davidsbond/kollect/internal/event"
	resource "github.com/davidsbond/kollect/proto/kollect/resource/event/v1"
)

type (
	// The EventHandler type is used to handle specific events published by kollect from a supported event bus.
	EventHandler struct {
		reader *event.Reader

		onResourceCreated ResourceCreatedHandler
		onResourceUpdated ResourceUpdatedHandler
		onResourceDeleted ResourceDeletedHandler
	}

	// The ResourceCreatedHandler type is a function that is invoked when the EventHandler consumes an event indicating
	// the creation/discovery of a new resource.
	ResourceCreatedHandler func(ctx context.Context, clusterID string, obj *unstructured.Unstructured) error

	// The ResourceUpdatedHandler type is a function that is invoked when the EventHandler consumes an event indicating
	// that an existing cluster resource has been modified.
	ResourceUpdatedHandler func(ctx context.Context, clusterID string, then, now *unstructured.Unstructured) error

	// The ResourceDeletedHandler type is a function that is invoked when the EventHandler consumes an event indicating
	// that an existing cluster resource was deleted.
	ResourceDeletedHandler func(ctx context.Context, clusterID, resourceUID string) error
)

// NewEventHandler returns a new instance of the EventHandler type that connects to the event bus described in the
// provided url.
func NewEventHandler(ctx context.Context, urlStr string) (*EventHandler, error) {
	reader, err := event.NewReader(ctx, urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to event bus: %w", err)
	}

	return &EventHandler{reader: reader}, nil
}

// OnResourceCreated sets up a ResourceCreatedHandler implementation to be invoked whenever an event that indicates a
// new resource is consumed.
func (eh *EventHandler) OnResourceCreated(fn ResourceCreatedHandler) {
	eh.onResourceCreated = fn
}

// OnResourceUpdated sets up a ResourceUpdatedHandler implementation to be invoked whenever an event that indicates an
// existing resource has been modified.
func (eh *EventHandler) OnResourceUpdated(fn ResourceUpdatedHandler) {
	eh.onResourceUpdated = fn
}

// OnResourceDeleted sets up a ResourceDeletedHandler implementation to be invoked whenever an event that indicates an
// existing resource has been deleted.
func (eh *EventHandler) OnResourceDeleted(fn ResourceDeletedHandler) {
	eh.onResourceDeleted = fn
}

// Handle inbound events, invoking any registered handler functions for their respective event types. This method
// blocks until the provided context is cancelled or one of the registered handler functions returns a non-nil
// error.
func (eh *EventHandler) Handle(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return eh.reader.Read(ctx, func(ctx context.Context, evt event.Event) error {
			switch payload := evt.Payload.(type) {
			case *resource.ResourceCreatedEvent:
				return eh.handleResourceCreatedEvent(ctx, payload)
			case *resource.ResourceUpdatedEvent:
				return eh.handleResourceUpdatedEvent(ctx, payload)
			case *resource.ResourceDeletedEvent:
				return eh.handleResourceDeletedEvent(ctx, payload)
			default:
				return nil
			}
		})
	})
	grp.Go(func() error {
		<-ctx.Done()
		return eh.reader.Close()
	})

	return grp.Wait()
}

func (eh *EventHandler) handleResourceCreatedEvent(ctx context.Context, payload *resource.ResourceCreatedEvent) error {
	if eh.onResourceCreated == nil {
		return nil
	}

	var obj unstructured.Unstructured
	if err := json.Unmarshal(payload.GetResource(), &obj); err != nil {
		return fmt.Errorf("failed to unmarshal resource %s: %w", payload.GetUid(), err)
	}

	return eh.onResourceCreated(ctx, payload.GetClusterId(), &obj)
}

func (eh *EventHandler) handleResourceUpdatedEvent(ctx context.Context, payload *resource.ResourceUpdatedEvent) error {
	if eh.onResourceUpdated == nil {
		return nil
	}

	var then, now unstructured.Unstructured
	if err := json.Unmarshal(payload.GetThen(), &then); err != nil {
		return fmt.Errorf("failed to unmarshal resource %s: %w", payload.GetUid(), err)
	}

	if err := json.Unmarshal(payload.GetNow(), &now); err != nil {
		return fmt.Errorf("failed to unmarshal resource %s: %w", payload.GetUid(), err)
	}

	return eh.onResourceUpdated(ctx, payload.GetClusterId(), &then, &now)
}

func (eh *EventHandler) handleResourceDeletedEvent(ctx context.Context, payload *resource.ResourceDeletedEvent) error {
	if eh.onResourceDeleted == nil {
		return nil
	}

	return eh.onResourceDeleted(ctx, payload.GetClusterId(), payload.GetUid())
}
