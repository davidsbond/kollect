# kollect

Kollect is a Kubernetes informer that monitors changes to all resources within your cluster and publishes those
changes as messages on your event bus of choice. Kollect aims to support as many event buses as possible. 

The following event buses are currently supported:

* [AWS SNS/SQS](https://aws.amazon.com/sns/)
* [GCP Pub/Sub](https://cloud.google.com/pubsub/docs/)
* [Azure Service Bus](https://azure.microsoft.com/en-us/services/service-bus/)
* [Apache Kafka](https://kafka.apache.org/)
* [NATS](https://nats.io/)
* [RabbitMQ](https://www.rabbitmq.com/)

Messages are encoded using [Protocol Buffers](https://developers.google.com/protocol-buffers), this is to make development 
of client implementations language-agnostic as well as minimising message size over the wire. You can view their definitions 
within the [proto](proto) directory.

## Getting started

Kollect can run both in and out-of cluster and requires a small number of command-line flags to operate. You can download
a binary from the [releases](https://github.com/davidsbond/kollect/releases) page, or pull the docker image for a release.

* `--cluster-id` (string): A unique identifier for the cluster that kollect is running in. This will allow clients to distinguish
between clusters when handling events.
* `--event-writer-url` (string): A URL that determines the event bus to use. Continue reading below for specifics on constructing
URLs for different event buses.
* `--kube-config` (string): The path to a kubeconfig file to use when running kollect outside a Kubernetes cluster. If you
intend to run kollect within a cluster, you can ignore this flag.
* `--namespace` (string): Configures kollect to only publish messages for resource changes within a particular namespace. Using 
this will mean that any cluster-scoped resources will not have updated published for them.
* `--wait-for-sync` (boolean): Configures kollect to wait for all informer caches to be synchronised before publishing any
messages. When unset, messages will be published as kollect builds the entire state of the cluster/namespace on startup.

## Event Bus URLs

Kollect configures its event writer via a URL whose scheme indicates the event bus to use. The underlying implementation
uses [gocloud.dev](https://gocloud.dev), their documentation should be used as a source of truth. Below are examples of
event bus URLs for each supported provider, and any additional environment variables required:

* AWS SNS/SQS

```
awssns:///arn:aws:sns:us-east-2:123456789012:mytopic?region=us-east-2
```

* GCP Pub/Sub

```
gcppubsub://projects/myproject/topics/mytopic
```

* Azure Service Bus

```
azuresb://mytopic
```

This event bus requires setting the `SERVICEBUS_CONNECTION_STRING` environment variable. This can be obtained from the 
[Azure portal](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-dotnet-how-to-use-topics-subscriptions#get-the-connection-string).

* Apache Kafka

```
kafka://my-topic
```

This event bus requires setting the `KAFKA_BROKERS` environment variable (which is a comma-delimited list of hosts, 
something like `1.2.3.4:9092,5.6.7.8:9092`).

* NATS

```
nats://example.mysubject
```

This event bus requires setting the `NATS_SERVER_URL` environment variable (which is something like `nats://nats.example.com`).

* RabbitMQ

```
rabbit://myexchange
```

This event bus requires setting the `RABBIT_SERVER_URL` environment variable (which is something like `amqp://guest:guest@localhost:5672/`).

## Monitoring

Kollect exposes a variety of endpoints on port `8081` to use for monitoring the application:

* `/__/metrics`: Serves [Prometheus](https://prometheus.io/) metrics.
* `/__/pprof`: Serves [pprof](https://github.com/google/pprof) endpoints for profiling.
* `/__/ready`: Serves an `HTTP OK` response when the application is considered ready.
* `/__/health`: Serves an `HTTP OK` response while the application is considered healthy.
