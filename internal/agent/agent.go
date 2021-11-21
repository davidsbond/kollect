// Package agent provides the implementation of the in-cluster agent that reacts to resource changes and sends data
// to a collector instance.
package agent

import (
	"context"
	"errors"
	"fmt"
	"path"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	"github.com/davidsbond/kollect/internal/event"
	"github.com/davidsbond/kollect/internal/telemetry"
	resource "github.com/davidsbond/kollect/proto/kollect/resource/event/v1"
)

type (
	// The Agent type is responsible for handling changes in resources within a cluster namespace and sending them
	// to a configured EventWriter.
	Agent struct {
		config Config

		// Flag used to prevent event writing until informer caches are synced.
		synced bool

		// Mutex used to ensure only a single handler is invoked at a time. For example, to prevent update events
		// being published before creation event. This mutex is wrapped with tracing spans so we can see any contention
		// in traces.
		handlerMux *telemetry.TracedMutex

		// Mutex used to get/set the synced flag across multiple goroutines.
		syncMux *sync.RWMutex
	}

	// The Config type describes configuration values that can be set for the Agent.
	Config struct {
		// The Namespace that resources will be collected for.
		Namespace string
		// The configuration for the cluster.
		ClusterClient dynamic.Interface
		// The resource types to send via the EventWriter.
		Resources []schema.GroupVersionResource
		// If true, no events are published until the initial informer caches are synced. This prevents events being
		// publishing describing the current state.
		WaitForCacheSync bool
		// The unique name for the cluster the agent is running in
		ClusterID string
		// The collector client used to consume resource changes.
		EventWriter EventWriter
	}

	// The EventWriter interface describes types that can publish events to an arbitrary event store.
	EventWriter interface {
		Write(ctx context.Context, evt event.Event) error
	}
)

// New returns a new instance of the Agent type with a set Config.
func New(config Config) *Agent {
	return &Agent{
		config:     config,
		handlerMux: telemetry.NewTracedMutex(),
		syncMux:    &sync.RWMutex{},
	}
}

var errCacheSyncFailed = errors.New("failed to sync cache")

const (
	keyKubernetesGroup     = "k8s.group"
	keyKubernetesVersion   = "k8s.version"
	keyKubernetesKind      = "k8s.kind"
	keyKubernetesName      = "k8s.name"
	keyKubernetesNamespace = "k8s.namespace"
	keyKubernetesUID       = "k8s.uid"
	keyClusterID           = "cluster.id"
)

// Run starts the agent, any detected changes in cluster resources will be sent to the configured EventWriter. Blocks until
// an error occurs or until the provided context.Context is cancelled.
func (a *Agent) Run(ctx context.Context) error {
	klog.Infoln("Agent starting")

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(a.config.ClusterClient, time.Minute*5, a.config.Namespace, nil)
	group, ctx := errgroup.WithContext(ctx)

	cacheSyncs := make([]cache.InformerSynced, len(a.config.Resources))
	for i, rs := range a.config.Resources {
		informer := factory.ForResource(rs).Informer()
		cacheSyncs[i] = informer.HasSynced

		handler := a.informerHandler(ctx, informer)
		group.Go(handler)
	}

	// Cache sync can be disabled if users want to build an initial state. Ideally this is only used to start with
	// then disabled.
	if !a.config.WaitForCacheSync {
		a.syncMux.Lock()
		a.synced = true
		a.syncMux.Unlock()

		return group.Wait()
	}

	// Return value from WaitForCacheSync is not assigned to a.synced directly within the Lock() and Unlock().
	// This is because handler functions are invoked while caches are syncing. If we lock around cache.WaitForCacheSync
	// the initial invocations of the add, update and delete handlers will still be waiting for the lock to be freed.
	// causing all those unwanted invocations to trigger events.
	synced := cache.WaitForCacheSync(ctx.Done(), cacheSyncs...)

	// Prevent any events from being written until the initial caches are synced. This prevents rewriting the entire
	// state of the cluster/namespace should the agent restart.
	a.syncMux.Lock()
	a.synced = synced
	a.syncMux.Unlock()

	if !a.isSynced() {
		return errCacheSyncFailed
	}

	return group.Wait()
}

// Ready returns true if all informer caches are synced.
func (a *Agent) Ready() bool {
	return a.isSynced()
}

func (a *Agent) informerHandler(ctx context.Context, informer cache.SharedIndexInformer) func() error {
	return func() error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc:    a.addHandler(ctx),
			UpdateFunc: a.updateHandler(ctx),
			DeleteFunc: a.deleteHandler(ctx),
		})
		err := informer.SetWatchErrorHandler(func(r *cache.Reflector, err error) {
			// If we don't have access to this resource, log and stop the informer so that we don't pollute the logs
			// doing this over and over again.
			klog.Errorln(err)
			cancel()
		})
		if err != nil {
			return err
		}

		go informer.Run(ctx.Done())
		<-ctx.Done()
		return nil
	}
}

func (a *Agent) addHandler(ctx context.Context) func(obj interface{}) {
	return func(obj interface{}) {
		if !a.isSynced() {
			return
		}

		err := telemetry.WithinSpan(ctx, "Resource.Create", func(ctx context.Context, span telemetry.Span) error {
			a.handlerMux.Lock(ctx)
			defer a.handlerMux.Unlock()

			item, ok := obj.(*unstructured.Unstructured)
			if !ok {
				klog.Errorln("item is not *unstructured.Unstructured")
				return nil
			}

			uid := string(item.GetUID())
			gvk := item.GroupVersionKind()

			span.SetAttributes(map[string]interface{}{
				keyKubernetesGroup:     gvk.Group,
				keyKubernetesVersion:   gvk.Version,
				keyKubernetesKind:      gvk.Kind,
				keyKubernetesNamespace: item.GetNamespace(),
				keyKubernetesName:      item.GetName(),
				keyKubernetesUID:       uid,
				keyClusterID:           a.config.ClusterID,
			})

			data, err := item.MarshalJSON()
			if err != nil {
				return fmt.Errorf("failed to marshal resource %s: %w", uid, err)
			}

			key := path.Join(a.config.ClusterID, uid)
			payload := &resource.ResourceCreatedEvent{
				Uid:       uid,
				Resource:  data,
				ClusterId: a.config.ClusterID,
			}

			evt := event.New(payload,
				event.WithKey(key),
				event.WithAppliesAt(item.GetCreationTimestamp().Time),
			)

			if err = a.config.EventWriter.Write(ctx, evt); err != nil {
				return fmt.Errorf("failed to publish event: %w", err)
			}

			resourceCreated.WithLabelValues(
				gvk.Group,
				gvk.Version,
				gvk.Kind,
				item.GetNamespace(),
			).Inc()
			return nil
		})
		if err != nil {
			klog.Errorln(err)
		}
	}
}

func (a *Agent) updateHandler(ctx context.Context) func(then, now interface{}) {
	return func(x, y interface{}) {
		if !a.isSynced() {
			return
		}

		err := telemetry.WithinSpan(ctx, "Resource.Update", func(ctx context.Context, span telemetry.Span) error {
			a.handlerMux.Lock(ctx)
			defer a.handlerMux.Unlock()

			then, ok := x.(*unstructured.Unstructured)
			if !ok {
				klog.Errorln("item is not *unstructured.Unstructured")
				return nil
			}

			now, ok := y.(*unstructured.Unstructured)
			if !ok {
				klog.Errorln("item is not *unstructured.Unstructured")
				return nil
			}

			uid := string(then.GetUID())
			gvk := then.GroupVersionKind()

			span.SetAttributes(map[string]interface{}{
				keyKubernetesGroup:     gvk.Group,
				keyKubernetesVersion:   gvk.Version,
				keyKubernetesKind:      gvk.Kind,
				keyKubernetesNamespace: then.GetNamespace(),
				keyKubernetesName:      then.GetName(),
				keyKubernetesUID:       then.GetUID(),
				keyClusterID:           a.config.ClusterID,
			})

			thenData, err := then.MarshalJSON()
			if err != nil {
				return fmt.Errorf("failed to marshal resource %s: %w", uid, err)
			}

			nowData, err := now.MarshalJSON()
			if err != nil {
				return fmt.Errorf("failed to marshal resource %s: %w", uid, err)
			}

			key := path.Join(a.config.ClusterID, uid)
			payload := &resource.ResourceUpdatedEvent{
				Uid:       uid,
				Then:      thenData,
				Now:       nowData,
				ClusterId: a.config.ClusterID,
			}

			evt := event.New(payload,
				event.WithKey(key),
				event.WithAppliesAt(time.Now()),
			)

			if err = a.config.EventWriter.Write(ctx, evt); err != nil {
				return fmt.Errorf("failed to publish event: %w", err)
			}

			resourceUpdated.WithLabelValues(
				gvk.Group,
				gvk.Version,
				gvk.Kind,
				then.GetNamespace(),
			).Inc()
			return nil
		})
		if err != nil {
			klog.Errorln(err)
		}
	}
}

func (a *Agent) deleteHandler(ctx context.Context) func(obj interface{}) {
	return func(obj interface{}) {
		if !a.isSynced() {
			return
		}

		err := telemetry.WithinSpan(ctx, "Resource.Delete", func(ctx context.Context, span telemetry.Span) error {
			a.handlerMux.Lock(ctx)
			defer a.handlerMux.Unlock()

			item, ok := obj.(*unstructured.Unstructured)
			if !ok {
				klog.Errorln("item is not *unstructured.Unstructured")
				return nil
			}

			gvk := item.GroupVersionKind()
			uid := string(item.GetUID())

			span.SetAttributes(map[string]interface{}{
				keyKubernetesGroup:     gvk.Group,
				keyKubernetesVersion:   gvk.Version,
				keyKubernetesKind:      gvk.Kind,
				keyKubernetesNamespace: item.GetNamespace(),
				keyKubernetesName:      item.GetName(),
				keyKubernetesUID:       uid,
				keyClusterID:           a.config.ClusterID,
			})

			key := path.Join(a.config.ClusterID, uid)
			payload := &resource.ResourceDeletedEvent{
				Uid:       uid,
				ClusterId: a.config.ClusterID,
			}

			evt := event.New(payload,
				event.WithKey(key),
				event.WithAppliesAt(item.GetCreationTimestamp().Time),
			)

			if err := a.config.EventWriter.Write(ctx, evt); err != nil {
				return fmt.Errorf("failed to publish event: %w", err)
			}

			resourceDeleted.WithLabelValues(
				gvk.Group,
				gvk.Version,
				gvk.Kind,
				item.GetNamespace(),
			).Inc()
			return nil
		})
		if err != nil {
			klog.Errorln(err)
		}
	}
}

func (a *Agent) isSynced() bool {
	a.syncMux.RLock()
	defer a.syncMux.RUnlock()
	return a.synced
}
