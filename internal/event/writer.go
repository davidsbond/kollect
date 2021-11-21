package event

import (
	"context"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/Shopify/sarama"
	"github.com/streadway/amqp"
	"gocloud.dev/pubsub"
	gcppubsub "google.golang.org/genproto/googleapis/pubsub/v1"

	"github.com/davidsbond/kollect/internal/telemetry"
)

type (
	// The Writer type is used to write events to a single topic.
	Writer struct {
		topic *pubsub.Topic
	}
)

// NewWriter creates a new instance of the Writer type that will write events to the configured
// event stream provider identified using the given URL.
func NewWriter(ctx context.Context, urlStr string) (*Writer, error) {
	topic, err := pubsub.OpenTopic(ctx, urlStr)
	return &Writer{topic: topic}, err
}

// Write an event to the stream.
func (w *Writer) Write(ctx context.Context, evt Event) error {
	return telemetry.WithinSpan(ctx, "Event.Write", func(ctx context.Context, span telemetry.Span) error {
		body, err := evt.marshal()
		if err != nil {
			return err
		}

		span.SetAttributes(map[string]interface{}{
			keyEventKey:  evt.Key,
			keyEventID:   evt.ID,
			keyEventType: evt.typeName(),
			keyEventSize: len(body),
		})

		err = w.topic.Send(ctx, &pubsub.Message{
			Body:       body,
			BeforeSend: producerKeyFunc(evt.Key),
			Metadata:   telemetry.Extract(ctx),
		})

		if err != nil {
			return err
		}

		eventsWritten.WithLabelValues(evt.Key, evt.typeName()).Inc()
		return nil
	})
}

// Close the connection to the event stream.
func (w *Writer) Close() error {
	const timeout = time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return w.topic.Shutdown(ctx)
}

// producerKeyFunc implements the pubsub.Message.BeforeSend function and is used to appropriately set the message key before
// publishing depending on the event bus being used. The asFunc is used to determine if the message about to be
// published is of a particular type. Documentation on the underlying types exposed by the As() functions can be
// seen in the gocloud documentation. https://gocloud.dev/concepts/as/#as
func producerKeyFunc(key string) func(asFunc func(interface{}) bool) error {
	if key == "" {
		return nil
	}

	// Return a function that modifies the contents of the message based on the type of message.
	return func(as func(interface{}) bool) error {
		pubsubMessage := &gcppubsub.PubsubMessage{}
		kafkaMessage := &sarama.ProducerMessage{}
		azureMessage := &servicebus.Message{}
		rabbitMessage := &amqp.Delivery{}

		switch {
		case as(&pubsubMessage):
			// If we're using a GCP pubsub message, set the ordering key field.
			pubsubMessage.OrderingKey = key
		case as(&kafkaMessage):
			// If we're using Kafka, set the partition key field.
			kafkaMessage.Key = sarama.StringEncoder(key)
		case as(&azureMessage):
			// If we're using Azure Service Bus, set the partition key property.
			azureMessage.SystemProperties.PartitionKey = &key
		case as(&rabbitMessage):
			// If we're using RabbitMQ, set the routing key.
			rabbitMessage.RoutingKey = key
		}

		return nil
	}
}
