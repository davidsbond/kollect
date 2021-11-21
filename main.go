// Package main contains the entrypoint to the agent application. The agent is responsible for monitoring changes in
// kubernetes cluster resources and publishing those changes as events onto a user-specified event bus.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"

	"github.com/davidsbond/kollect/internal/agent"
	"github.com/davidsbond/kollect/internal/event"
	"github.com/davidsbond/kollect/internal/kubernetes"
)

var version string

func main() {
	var (
		eventWriterURL string
		kubeConfig     string
		namespace      string
		waitForSync    bool
		clusterID      string
	)

	closer := func(c io.Closer) {
		if err := c.Close(); err != nil {
			klog.Errorf("failed to close %T: %v\n", c, err)
		}
	}

	run := func(ctx context.Context) error {
		eventWriter, err := event.NewWriter(ctx, eventWriterURL)
		if err != nil {
			return fmt.Errorf("failed to connect to event bus: %w", err)
		}
		defer closer(eventWriter)

		k8sConfig, err := kubernetes.Config(kubeConfig)
		if err != nil {
			return fmt.Errorf("failed to create k8s config: %w", err)
		}

		cnf := agent.Config{
			EventWriter:      eventWriter,
			Namespace:        namespace,
			WaitForCacheSync: waitForSync,
			ClusterID:        clusterID,
		}

		cnf.Resources, err = kubernetes.GetResourcesWithVerbs(k8sConfig, []string{"get", "list", "watch"})
		if err != nil {
			return fmt.Errorf("failed to list k8s resources: %w", err)
		}

		cnf.ClusterClient, err = dynamic.NewForConfig(k8sConfig)
		if err != nil {
			return fmt.Errorf("failed to create dynamic k8s client: %w", err)
		}

		ag := agent.New(cnf)
		grp, ctx := errgroup.WithContext(ctx)
		grp.Go(func() error {
			return ag.Run(ctx)
		})

		grp.Go(func() error {
			mux := http.NewServeMux()
			svr := &http.Server{Addr: ":8081", Handler: mux}

			mux.Handle("/__/metrics", promhttp.Handler())
			mux.HandleFunc("/__/pprof/profile", pprof.Profile)
			mux.HandleFunc("/__/pprof/trace", pprof.Trace)
			mux.HandleFunc("/__/pprof/cmdline", pprof.Cmdline)
			mux.HandleFunc("/__/pprof/symbol", pprof.Symbol)

			mux.HandleFunc("/__/ready", func(w http.ResponseWriter, r *http.Request) {
				if !ag.Ready() {
					w.WriteHeader(http.StatusPreconditionFailed)
				}
			})

			mux.HandleFunc("/__/health", func(w http.ResponseWriter, r *http.Request) {})

			return svr.ListenAndServe()
		})

		return grp.Wait()
	}

	cmd := &cobra.Command{
		Use:     "kollect",
		Short:   "Publish changes in your Kubernetes resources as events on your choice of event bus",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			if err := run(ctx); err != nil {
				klog.Exitln(err)
			}
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&namespace, "namespace", v1.NamespaceAll, "Specifies the namespace that the agent will monitor resources in, defaults to all")
	flags.StringVar(&eventWriterURL, "event-writer-url", "", "URL of the event bus to send resource events to, see documentation for possible values")
	flags.StringVar(&kubeConfig, "kube-config", "", "Location of the kubeconfig file to use for authentication. In-cluster config used if blank")
	flags.BoolVar(&waitForSync, "wait-for-sync", false, "If set, no events will be published until the caches are synced. When false, events will be published for the entire cluster state on start")
	flags.StringVar(&clusterID, "cluster-id", "", "The unique identifier for the cluster the agent is running in")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := cmd.ExecuteContext(ctx); err != nil {
		klog.Exitln(err)
	}
}
