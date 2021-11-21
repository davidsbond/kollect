package kollect_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/davidsbond/kollect/internal/event"
	"github.com/davidsbond/kollect/pkg/kollect"
	resource "github.com/davidsbond/kollect/proto/kollect/resource/event/v1"
)

func TestEventHandler_Handle(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name              string
		Event             event.Event
		ExpectedClusterID string
		ExpectedResource  *unstructured.Unstructured
		ExpectsError      bool
	}{
		{
			Name: "It should handle a resource created event",
			Event: event.New(&resource.ResourceCreatedEvent{
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
			}),
			ExpectedClusterID: "test",
			ExpectedResource: &unstructured.Unstructured{
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
		},
		{
			Name: "It should handle a resource updated event",
			Event: event.New(&resource.ResourceUpdatedEvent{
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
						},
					},
				}),
				ClusterId: "test",
			}),
			ExpectedClusterID: "test",
			ExpectedResource: &unstructured.Unstructured{
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
		},
		{
			Name: "It should handle a resource deleted event",
			Event: event.New(&resource.ResourceDeletedEvent{
				Uid:       "test",
				ClusterId: "test",
			}),
			ExpectedClusterID: "test",
			ExpectedResource: &unstructured.Unstructured{
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
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			writer, err := event.NewWriter(ctx, "mem://test")
			require.NoError(t, err)

			handler, err := kollect.NewEventHandler(ctx, "mem://test")
			require.NoError(t, err)
			require.NoError(t, writer.Write(ctx, tc.Event))

			handler.OnResourceCreated(func(ctx context.Context, clusterID string, obj *unstructured.Unstructured) error {
				cancel()

				assert.EqualValues(t, tc.ExpectedClusterID, clusterID)
				assert.EqualValues(t, tc.ExpectedResource, obj)
				return nil
			})

			handler.OnResourceUpdated(func(ctx context.Context, clusterID string, then, now *unstructured.Unstructured) error {
				cancel()

				assert.EqualValues(t, tc.ExpectedClusterID, clusterID)
				assert.EqualValues(t, tc.ExpectedResource, then)
				assert.EqualValues(t, tc.ExpectedResource, now)
				return nil
			})

			handler.OnResourceDeleted(func(ctx context.Context, clusterID, resourceUID string) error {
				cancel()

				assert.EqualValues(t, tc.ExpectedClusterID, clusterID)
				assert.EqualValues(t, tc.ExpectedResource.GetUID(), resourceUID)
				return nil
			})

			err = handler.Handle(ctx)
			if tc.ExpectsError {
				assert.Error(t, err)
				assert.NotEqual(t, context.Canceled, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func mustMarshal(t *testing.T, obj *unstructured.Unstructured) []byte {
	t.Helper()

	data, err := obj.MarshalJSON()
	require.NoError(t, err)

	return data
}
