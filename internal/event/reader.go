package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/Shopify/sarama"
	"github.com/streadway/amqp"
	"gocloud.dev/pubsub"
	gcppubsub "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type (
	// The Reader type is used to handle inbound events from a single topic.
	Reader struct {
		subscription *pubsub.Subscription
	}
)

// NewReader creates a new instance of the Reader type that will read events from the configured
// event stream provider identified using the given URL.
func NewReader(ctx context.Context, urlStr string) (*Reader, error) {
	subscription, err := pubsub.OpenSubscription(ctx, urlStr)
	if err != nil {
		return nil, err
	}

	return &Reader{subscription: subscription}, nil
}

// Read events from the stream, invoking fn for each inbound event. This method will block until fn returns
// an error when messages are not nackable or the provided context is cancelled.
func (r *Reader) Read(ctx context.Context, fn Handler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := r.subscription.Receive(ctx)
			if err != nil {
				return err
			}

			evt, err := unmarshal(msg.Body)
			switch {
			case errors.Is(err, protoregistry.NotFound):
				continue
			case err != nil:
				nack(msg)
				return err
			}

			evt.Key = consumerKey(msg)

			err = fn(ctx, evt)
			switch {
			case err != nil:
				nack(msg)
				return fmt.Errorf("failed to handle event %s: %w", evt.ID, err)
			default:
				msg.Ack()
				eventsRead.WithLabelValues(evt.Key, evt.typeName()).Inc()
				return nil
			}
		}
	}
}

func nack(msg *pubsub.Message) {
	if msg.Nackable() {
		msg.Nack()
	}
}

// Close the connection to the event stream.
func (r *Reader) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return r.subscription.Shutdown(ctx)
}

func consumerKey(msg *pubsub.Message) string {
	pubsubMessage := &gcppubsub.PubsubMessage{}
	kafkaMessage := &sarama.ConsumerMessage{}
	azureMessage := &servicebus.Message{}
	rabbitMessage := &amqp.Delivery{}

	switch {
	case msg.As(&pubsubMessage):
		return pubsubMessage.OrderingKey
	case msg.As(&kafkaMessage):
		return string(kafkaMessage.Key)
	case msg.As(&azureMessage):
		if azureMessage.SystemProperties.PartitionKey != nil {
			return *azureMessage.SystemProperties.PartitionKey
		}
	case msg.As(&rabbitMessage):
		return rabbitMessage.RoutingKey
	}

	return ""
}
