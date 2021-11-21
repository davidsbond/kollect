package kollect

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/davidsbond/kollect/internal/event"
	resource "github.com/davidsbond/kollect/proto/kollect/resource/event/v1"
)

type (
	EventHandler struct {
		reader EventReader

		onResourceCreated ResourceCreatedHandler
		onResourceUpdated ResourceUpdatedHandler
		onResourceDeleted ResourceDeletedHandler
	}

	ResourceCreatedHandler func(ctx context.Context, obj *unstructured.Unstructured) error
	ResourceUpdatedHandler func(ctx context.Context, then, now *unstructured.Unstructured) error
	ResourceDeletedHandler func(ctx context.Context, clusterID, resourceUID string) error

	EventReader interface {
		Read(ctx context.Context, fn event.Handler) error
	}
)

func NewEventHandler(ctx context.Context, urlStr string) (*EventHandler, error) {
	reader, err := event.NewReader(ctx, urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to event bus: %w", err)
	}

	return &EventHandler{reader: reader}, nil
}

func (eh *EventHandler) OnResourceCreated(fn ResourceCreatedHandler) {
	eh.onResourceCreated = fn
}

func (eh *EventHandler) OnResourceUpdated(fn ResourceUpdatedHandler) {
	eh.onResourceUpdated = fn
}

func (eh *EventHandler) OnResourceDeleted(fn ResourceDeletedHandler) {
	eh.onResourceDeleted = fn
}

func (eh *EventHandler) Handle(ctx context.Context) error {
	return eh.reader.Read(ctx, func(ctx context.Context, evt event.Event) error {
		switch payload := evt.Payload.(type) {
		case *resource.ResourceCreatedEvent:
			return eh.handleResourceCreatedEvent(ctx, payload)
		case *resource.ResourceUpdatedEvent:
			return eh.handleResourceUpdatedEvent(ctx, payload)
		case *resource.ResourceDeletedEvent:
			return eh.handleResourceDeletedEvent(ctx, payload)
		default:
			return event.Ignore
		}
	})
}

func (eh *EventHandler) handleResourceCreatedEvent(ctx context.Context, payload *resource.ResourceCreatedEvent) error {
	if eh.onResourceCreated == nil {
		return event.Ignore
	}

	var obj unstructured.Unstructured
	if err := json.Unmarshal(payload.GetResource(), &obj); err != nil {
		return fmt.Errorf("failed to unmarshal resource %s: %w", payload.GetUid(), err)
	}

	return eh.onResourceCreated(ctx, &obj)
}

func (eh *EventHandler) handleResourceUpdatedEvent(ctx context.Context, payload *resource.ResourceUpdatedEvent) error {
	if eh.onResourceUpdated == nil {
		return event.Ignore
	}

	var then, now unstructured.Unstructured
	if err := json.Unmarshal(payload.GetThen(), &then); err != nil {
		return fmt.Errorf("failed to unmarshal resource %s: %w", payload.GetUid(), err)
	}

	if err := json.Unmarshal(payload.GetNow(), &now); err != nil {
		return fmt.Errorf("failed to unmarshal resource %s: %w", payload.GetUid(), err)
	}

	return eh.onResourceUpdated(ctx, &then, &now)
}

func (eh *EventHandler) handleResourceDeletedEvent(ctx context.Context, payload *resource.ResourceDeletedEvent) error {
	if eh.onResourceDeleted == nil {
		return event.Ignore
	}

	return eh.onResourceDeleted(ctx, payload.GetClusterId(), payload.GetUid())
}
