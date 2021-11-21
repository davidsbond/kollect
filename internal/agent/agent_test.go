package agent_test

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/fake"

	"github.com/davidsbond/kollect/internal/agent"
	resource "github.com/davidsbond/kollect/proto/kollect/resource/event/v1"
)

var (
	k8sClient       *fake.FakeDynamicClient
	testEventWriter *MockEventWriter
)

func TestMain(m *testing.M) {
	k8sClient = fake.NewSimpleDynamicClientWithCustomListKinds(runtime.NewScheme(), map[schema.GroupVersionResource]string{
		schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}: "UnstructuredList",
	})

	testEventWriter = &MockEventWriter{emitted: make(chan bool, 1)}

	cnf := agent.Config{
		Namespace:     "namespace",
		EventWriter:   testEventWriter,
		ClusterClient: k8sClient,
		ClusterID:     "test",
		Resources: []schema.GroupVersionResource{
			{Group: "apps", Version: "v1", Resource: "deployments"},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := agent.New(cnf).Run(ctx); err != nil {
			log.Fatalln(err)
		}
	}()

	<-time.After(time.Second)
	os.Exit(m.Run())
}

func TestAgent_Run(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		Name            string
		Resource        *unstructured.Unstructured
		ExpectedPayload proto.Message
		Action          func(t *testing.T, client dynamic.ResourceInterface, resource *unstructured.Unstructured)
	}{
		{
			Name: "It should emit created resources",
			Action: func(t *testing.T, client dynamic.ResourceInterface, resource *unstructured.Unstructured) {
				_, err := client.Create(ctx, resource, metav1.CreateOptions{
					TypeMeta: metav1.TypeMeta{
						Kind:       resource.GetKind(),
						APIVersion: resource.GetAPIVersion(),
					},
				})

				assert.NoError(t, err)
			},
			Resource: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name":      "example",
						"namespace": "namespace",
						"uid":       "test",
					},
				},
			},
			ExpectedPayload: &resource.ResourceCreatedEvent{
				Uid: "test",
				Resource: mustMarshal(t, &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "apps/v1",
						"kind":       "Deployment",
						"metadata": map[string]interface{}{
							"name":      "example",
							"namespace": "namespace",
							"uid":       "test",
						},
					},
				}),
				ClusterId: "test",
			},
		},
		{
			Name: "It should emit updated resources",
			Action: func(t *testing.T, client dynamic.ResourceInterface, resource *unstructured.Unstructured) {
				_, err := client.Update(ctx, resource, metav1.UpdateOptions{
					TypeMeta: metav1.TypeMeta{
						Kind:       resource.GetKind(),
						APIVersion: resource.GetAPIVersion(),
					},
				})

				assert.NoError(t, err)
			},
			Resource: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name":      "example",
						"namespace": "namespace",
						"uid":       "test",
						"labels": map[string]interface{}{
							"test-label": "label-value",
						},
					},
				},
			},
			ExpectedPayload: &resource.ResourceUpdatedEvent{
				Uid: "test",
				Then: mustMarshal(t, &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "apps/v1",
						"kind":       "Deployment",
						"metadata": map[string]interface{}{
							"name":      "example",
							"namespace": "namespace",
							"uid":       "test",
						},
					},
				}),
				Now: mustMarshal(t, &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "apps/v1",
						"kind":       "Deployment",
						"metadata": map[string]interface{}{
							"name":      "example",
							"namespace": "namespace",
							"uid":       "test",
							"labels": map[string]interface{}{
								"test-label": "label-value",
							},
						},
					},
				}),
				ClusterId: "test",
			},
		},
		{
			Name: "It should collect deleted resources",
			Action: func(t *testing.T, client dynamic.ResourceInterface, resource *unstructured.Unstructured) {
				assert.NoError(t, client.Delete(ctx, resource.GetName(), metav1.DeleteOptions{
					TypeMeta: metav1.TypeMeta{
						Kind:       resource.GetKind(),
						APIVersion: resource.GetAPIVersion(),
					},
				}))
			},
			ExpectedPayload: &resource.ResourceDeletedEvent{
				Uid:       "test",
				ClusterId: "test",
			},
			Resource: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name":      "example",
						"namespace": "namespace",
						"uid":       "test",
						"labels": map[string]interface{}{
							"test-label": "label-value",
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			gvr := schema.GroupVersionResource{
				Group:    tc.Resource.GroupVersionKind().Group,
				Version:  tc.Resource.GroupVersionKind().Version,
				Resource: strings.ToLower(tc.Resource.GetKind()) + "s",
			}

			client := k8sClient.Resource(gvr).Namespace(tc.Resource.GetNamespace())
			tc.Action(t, client, tc.Resource)

			testEventWriter.Wait()

			assert.True(t, proto.Equal(tc.ExpectedPayload, testEventWriter.event.Payload))
		})
	}
}

func mustMarshal(t *testing.T, obj *unstructured.Unstructured) []byte {
	t.Helper()

	data, err := obj.MarshalJSON()
	require.NoError(t, err)

	return data
}
