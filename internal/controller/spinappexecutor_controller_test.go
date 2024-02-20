package controller

import (
	"context"
	"testing"
	"time"

	spinv1 "github.com/spinkube/spin-operator/api/v1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

func TestSpinAppExecutorReconcilerStartupShutdown(t *testing.T) {
	t.Parallel()

	envTest := SetupEnvTest(t)

	mgr, err := ctrl.NewManager(envTest.cfg, manager.Options{
		Metrics: metricsserver.Options{BindAddress: "0"},
	})
	require.NoError(t, err)

	reconciler := &SpinAppExecutorReconciler{
		Client:   envTest.k8sClient,
		Scheme:   scheme.Scheme,
		Recorder: mgr.GetEventRecorderFor("spinappexecutor-controller"),
	}

	require.NoError(t, reconciler.SetupWithManager(mgr))

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	require.NoError(t, mgr.Start(ctx))
}

func TestDeleteSpinAppExecutor(t *testing.T) {
	t.Parallel()

	envTest := SetupEnvTest(t)

	mgr, err := ctrl.NewManager(envTest.cfg, manager.Options{
		Metrics: metricsserver.Options{BindAddress: "0"},
	})
	require.NoError(t, err)

	reconciler := &SpinAppExecutorReconciler{
		Client:   envTest.k8sClient,
		Scheme:   scheme.Scheme,
		Recorder: mgr.GetEventRecorderFor("spinappexecutor-controller"),
	}

	require.NoError(t, reconciler.SetupWithManager(mgr))

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	go func() {
		require.NoError(t, mgr.Start(ctx))
	}()

	// Create SpinAppExecutor
	err = envTest.k8sClient.Create(ctx, &spinv1.SpinAppExecutor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerd-shim-spin",
			Namespace: "default",
		},
	})
	require.NoError(t, err)

	// Create dependent SpinApp
	err = envTest.k8sClient.Create(ctx, &spinv1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spinapp",
			Namespace: "default",
		},
		Spec: spinv1.SpinAppSpec{
			Executor: "containerd-shim-spin",
			Replicas: 1,
			Image:    "ghcr.io/deislabs/containerd-wasm-shims/examples/spin-rust-hello:v0.10.0",
		},
	})
	require.NoError(t, err)

	// Attempt to delete SpinAppExecutor
	innerCtx, innerCancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer innerCancelFunc()
	err = envTest.k8sClient.Delete(innerCtx, &spinv1.SpinAppExecutor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerd-shim-spin",
			Namespace: "default",
		},
	})
	require.NoError(t, err)

	// Get SpinAppExecutor and verify the deletion timestamp is set
	executor := &spinv1.SpinAppExecutor{}
	err = envTest.k8sClient.Get(ctx, client.ObjectKey{
		Name:      "containerd-shim-spin",
		Namespace: "default",
	}, executor)
	require.NoError(t, err)
	require.False(t, executor.DeletionTimestamp.IsZero())

	// Delete dependent SpinApp
	err = envTest.k8sClient.Delete(ctx, &spinv1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spinapp",
			Namespace: "default",
		},
	})
	require.NoError(t, err)

	// TODO: Verify that the SpinAppExecutor is deleted
}
