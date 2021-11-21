// Package kubernetes provides functions for interacting with the kubernetes API.
package kubernetes

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"

	// Ensures different auth types works for the k8s client.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Config returns a rest.Config instance configured to use the kube config file located at the given path. If the
// kubeConfig parameter is blank, an in-cluster configuration is assumed.
func Config(kubeConfig string) (*rest.Config, error) {
	var err error
	var clusterConfig *rest.Config
	if kubeConfig != "" {
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	} else {
		clusterConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	return clusterConfig, nil
}

// GetResourcesWithVerbs returns all resources available in the cluster that can be accessed by all the verbs within
// the slice.
func GetResourcesWithVerbs(config *rest.Config, verbs []string) ([]schema.GroupVersionResource, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	lists, err := client.Discovery().ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	resources := make([]schema.GroupVersionResource, 0)
	for _, list := range lists {
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			continue
		}

		for _, resource := range list.APIResources {
			if !containsAll(resource.Verbs, verbs) {
				continue
			}

			resources = append(resources, schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: resource.Name,
			})
		}
	}

	return resources, nil
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
