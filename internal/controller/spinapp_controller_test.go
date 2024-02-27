package controller

import (
	"context"
	"fmt"
	"path/filepath"
	goruntime "runtime"
	"testing"
	"time"

	spinv1 "github.com/spinkube/spin-operator/api/v1"
	"github.com/stretchr/testify/require"
	runtime "k8s.io/apimachinery/pkg/runtime"
	clientscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

type envTestState struct {
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	scheme    *runtime.Scheme
}

// SetupEnvTest will start a fake kubernetes and client for use when testing
// reconciliation loops that require a kubernetes api.
func SetupEnvTest(t *testing.T) *envTestState {
	t.Helper()

	testEnv := &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without calling the makefile target test. If not informed it will look for the
		// default path defined in controller-runtime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("..", "..", "bin", "k8s",
			fmt.Sprintf("1.28.3-%s-%s", goruntime.GOOS, goruntime.GOARCH)),
	}

	cfg, err := testEnv.Start()
	if err != nil {
		t.Skipf("envtest unavailable: %v", err)
	}

	require.NoError(t, err)
	require.NotNil(t, cfg)

	scheme := runtime.NewScheme()
	clientscheme.AddToScheme(scheme)

	err = spinv1.AddToScheme(scheme)
	require.NoError(t, err)

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme})
	require.NoError(t, err)
	require.NotNil(t, k8sClient)

	t.Cleanup(func() {
		err := testEnv.Stop()
		require.NoError(t, err)
	})

	return &envTestState{
		cfg:       cfg,
		k8sClient: k8sClient,
		testEnv:   testEnv,
		scheme:    scheme,
	}
}

func TestReconcile_Integration_StartupShutdown(t *testing.T) {
	t.Parallel()

	envTest := SetupEnvTest(t)

	mgr, err := ctrl.NewManager(envTest.cfg, manager.Options{
		Metrics: metricsserver.Options{BindAddress: "0"},
		Scheme:  envTest.scheme,
	})

	require.NoError(t, err)

	ctrlr := &SpinAppReconciler{
		Client:   envTest.k8sClient,
		Scheme:   envTest.scheme,
		Recorder: mgr.GetEventRecorderFor("spinapp-reconciler"),
	}

	require.NoError(t, ctrlr.SetupWithManager(mgr))

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	require.NoError(t, mgr.Start(ctx))
}

func TestConstructDeployment_MinimalApp(t *testing.T) {
	t.Parallel()

	app := minimalSpinApp()

	cfg := &spinv1.ExecutorDeploymentConfig{
		RuntimeClassName: "bananarama",
	}
	deployment, err := constructDeployment(context.Background(), app, cfg, nil)
	require.NoError(t, err)
	require.NotNil(t, deployment)

	require.Equal(t, ptr(int32(1)), deployment.Spec.Replicas)
	require.Len(t, deployment.Spec.Template.Spec.Containers, 1)
	require.Equal(t, app.Spec.Image, deployment.Spec.Template.Spec.Containers[0].Image)
	require.Equal(t, ptr("bananarama"), deployment.Spec.Template.Spec.RuntimeClassName)
}
