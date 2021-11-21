// Package main contains the entrypoint to the agent application. The agent is responsible for monitoring changes in
// kubernetes cluster resources and publishing those changes as events onto a user-specified event bus.
package main

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/davidsbond/kollect/internal/agent"
	"github.com/davidsbond/kollect/internal/app"
	"github.com/davidsbond/kollect/internal/closers"
	"github.com/davidsbond/kollect/internal/event"
	"github.com/davidsbond/kollect/internal/flag"
	"github.com/davidsbond/kollect/internal/health"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "namespace",
				Usage:       "Specifies the namespace that the agent will monitor resources in, defaults to all.",
				Destination: &namespace,
				EnvVar:      "NAMESPACE",
				Value:       v1.NamespaceAll,
			},
			&flag.String{
				Name:        "event-writer-url",
				Usage:       "URL of the event bus to send resource events to",
				Destination: &eventWriterURL,
				EnvVar:      "EVENT_WRITER_URL",
				Required:    true,
			},
			&flag.String{
				Name:        "kube-config",
				Usage:       "Location of the kubeconfig file to use for authentication. In-cluster config used if blank.",
				Destination: &kubeConfig,
				EnvVar:      "KUBECONFIG",
			},
			&flag.Boolean{
				Name:        "wait-for-sync",
				Usage:       "If true, no events will be published until the caches are synced. When false, events will be published for the entire cluster state on start.",
				Destination: &waitForSync,
				EnvVar:      "WAIT_FOR_SYNC",
				Value:       true,
			},
			&flag.String{
				Name:        "cluster-id",
				Usage:       "The unique identifier for the cluster the agent is running in",
				Destination: &clusterID,
				EnvVar:      "CLUSTER_ID",
				Required:    true,
			},
		),
	)

	if err := a.Run(); err != nil {
		log.Fatalln(err)
	}
}

var (
	eventWriterURL string
	kubeConfig     string
	namespace      string
	waitForSync    bool
	clusterID      string
)

func run(ctx context.Context) error {
	eventWriter, err := event.NewWriter(ctx, eventWriterURL)
	if err != nil {
		return err
	}
	defer closers.Close(eventWriter)

	cnf := agent.Config{
		EventWriter:      eventWriter,
		Namespace:        namespace,
		WaitForCacheSync: waitForSync,
		ClusterID:        clusterID,
	}

	var clusterConfig *rest.Config
	if kubeConfig != "" {
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	} else {
		clusterConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return err
	}

	health.AddHealthCheck(clusterConfig.Host, health.CheckKubernetesAPI(clusterConfig))

	cl, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return err
	}

	lists, err := cl.Discovery().ServerPreferredResources()
	if err != nil {
		return err
	}

	for _, list := range lists {
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			log.Println("failed to parse GroupVersion", list.GroupVersion)
			continue
		}

		for _, resource := range list.APIResources {
			// Attempting to watch cluster scoped resources while using a namespaced informer will return
			// an error. So we check if the configured namespace is explicitly set from the default and that
			// the resource is cluster scoped.
			if namespace != v1.NamespaceAll && !resource.Namespaced {
				continue
			}

			// For informers to work correctly, we have to ensure the resource can be retrieved, watched and
			// listened to.
			if !containsAll(resource.Verbs, []string{"get", "list", "watch"}) {
				continue
			}

			cnf.Resources = append(cnf.Resources, schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: resource.Name,
			})
		}
	}

	cnf.ClusterClient, err = dynamic.NewForConfig(clusterConfig)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	a := agent.New(cnf)
	health.AddReadyCheck(a.Ready)

	return a.Run(ctx)
}

func containsAll(arr, values []string) bool {
	for _, value := range values {
		contains := false

		for _, elem := range arr {
			contains = value == elem
			if contains {
				break
			}
		}

		if !contains {
			return false
		}
	}

	return true
}
