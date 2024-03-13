package webhook

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"path/filepath"
	goruntime "runtime"
	"testing"
	"time"

	spinv1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/internal/constants"
	"github.com/stretchr/testify/require"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type envTestState struct {
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	scheme    *runtime.Scheme
}

// setupEnvTest will start a fake kubernetes and client for use when testing
// webhooks that require a kubernetes api.
func setupEnvTest(t *testing.T) *envTestState {
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

		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: []string{filepath.Join("..", "..", "config", "webhook")},
		},
	}

	cfg, err := testEnv.Start()
	if err != nil {
		t.Skipf("envtest unavailable: %v", err)
	}

	require.NoError(t, err)
	require.NotNil(t, cfg)

	scheme := runtime.NewScheme()
	require.NoError(t, clientscheme.AddToScheme(scheme))
	err = spinv1alpha1.AddToScheme(scheme)
	require.NoError(t, err)

	err = admissionv1.AddToScheme(scheme)
	require.NoError(t, err)

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme})
	require.NoError(t, err)
	require.NotNil(t, k8sClient)

	return &envTestState{
		cfg:       cfg,
		k8sClient: k8sClient,
		testEnv:   testEnv,
		scheme:    scheme,
	}
}

func startWebhookServer(t *testing.T, envtest *envTestState) {
	t.Helper()

	// start webhook server using Manager
	webhookInstallOptions := &envtest.testEnv.WebhookInstallOptions
	mgr, err := ctrl.NewManager(envtest.cfg, ctrl.Options{
		Scheme: envtest.scheme,
		WebhookServer: webhook.NewServer(webhook.Options{
			Host:    webhookInstallOptions.LocalServingHost,
			Port:    webhookInstallOptions.LocalServingPort,
			CertDir: webhookInstallOptions.LocalServingCertDir,
		}),
		LeaderElection: false,
		Metrics:        metricsserver.Options{BindAddress: "0"},
	})
	require.NoError(t, err)

	err = SetupSpinAppWebhookWithManager(mgr)
	require.NoError(t, err)

	err = SetupSpinAppExecutorWebhookWithManager(mgr)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err = mgr.Start(ctx)
		require.NoError(t, err)
	}()

	// wait for the webhook server to get ready
	dialer := &net.Dialer{Timeout: time.Second}
	addrPort := fmt.Sprintf("%s:%d", webhookInstallOptions.LocalServingHost, webhookInstallOptions.LocalServingPort)
	require.Eventually(t, func() bool {
		// nolint:gosec
		conn, err := tls.DialWithDialer(dialer, "tcp", addrPort, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return false
		}
		err = conn.Close()
		return err == nil
	}, 10*time.Second, 2*time.Second)

	t.Cleanup(func() {
		err := envtest.testEnv.Stop()
		require.NoError(t, err)
	})

	// As per https://github.com/kubernetes-sigs/controller-runtime/issues/1571 to avoid leaking kube-apiserver and etcd
	t.Cleanup(func() {
		cancel()
	})
}

func TestCreateSpinAppWithNoExecutor(t *testing.T) {
	t.Parallel()

	envtest := setupEnvTest(t)
	startWebhookServer(t, envtest)

	err := envtest.k8sClient.Create(context.Background(), &spinv1alpha1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spinapp",
			Namespace: "default",
		},
		Spec: spinv1alpha1.SpinAppSpec{
			Image:    "ghcr.io/spinkube/containerd-shim-spin/examples/spin-rust-hello:v0.13.0",
			Replicas: 2,
		},
	})
	require.EqualError(t, err, "admission webhook \"vspinapp.kb.io\" denied the request: SpinApp.core.spinoperator.dev \"spinapp\" is invalid:"+
		" spec.executor: Invalid value: \"\": executor must be set, likely no default executor was set because you have no executors installed")
}

func TestCreateSpinAppWithSingleExecutor(t *testing.T) {
	t.Parallel()

	envtest := setupEnvTest(t)
	startWebhookServer(t, envtest)

	err := envtest.k8sClient.Create(context.Background(), &spinv1alpha1.SpinAppExecutor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cyclotron",
			Namespace: "default",
		},
	})
	require.NoError(t, err)

	err = envtest.k8sClient.Create(context.Background(), &spinv1alpha1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spinapp",
			Namespace: "default",
		},
		Spec: spinv1alpha1.SpinAppSpec{
			Image:    "ghcr.io/spinkube/containerd-shim-spin/examples/spin-rust-hello:v0.13.0",
			Replicas: 2,
		},
	})
	require.NoError(t, err)

	spinapp := &spinv1alpha1.SpinApp{}
	err = envtest.k8sClient.Get(context.Background(), client.ObjectKey{
		Name:      "spinapp",
		Namespace: "default",
	}, spinapp)
	require.NoError(t, err)
	require.Equal(t, constants.CyclotronExecutor, spinapp.Spec.Executor)
	require.Equal(t, int32(2), spinapp.Spec.Replicas)
}

func TestCreateSpinAppWithMultipleExecutors(t *testing.T) {
	t.Parallel()

	envtest := setupEnvTest(t)
	startWebhookServer(t, envtest)

	err := envtest.k8sClient.Create(context.Background(), &spinv1alpha1.SpinAppExecutor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerd-shim-spin",
			Namespace: "default",
		},
	})
	require.NoError(t, err)

	err = envtest.k8sClient.Create(context.Background(), &spinv1alpha1.SpinAppExecutor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cyclotron",
			Namespace: "default",
		},
	})
	require.NoError(t, err)

	err = envtest.k8sClient.Create(context.Background(), &spinv1alpha1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spinapp",
			Namespace: "default",
		},
		Spec: spinv1alpha1.SpinAppSpec{
			Image:    "ghcr.io/spinkube/containerd-shim-spin/examples/spin-rust-hello:v0.13.0",
			Replicas: 2,
		},
	})
	require.NoError(t, err)

	spinapp := &spinv1alpha1.SpinApp{}
	err = envtest.k8sClient.Get(context.Background(), client.ObjectKey{
		Name:      "spinapp",
		Namespace: "default",
	}, spinapp)
	require.NoError(t, err)
	// Correct based on alphabetical order
	require.Equal(t, constants.ContainerDShimSpinExecutor, spinapp.Spec.Executor)
	require.Equal(t, int32(2), spinapp.Spec.Replicas)
}

func TestCreateInvalidSpinApp(t *testing.T) {
	t.Parallel()

	envtest := setupEnvTest(t)
	startWebhookServer(t, envtest)

	err := envtest.k8sClient.Create(context.Background(), &spinv1alpha1.SpinAppExecutor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerd-shim-spin",
			Namespace: "default",
		},
	})
	require.NoError(t, err)

	err = envtest.k8sClient.Create(context.Background(), &spinv1alpha1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spinapp",
			Namespace: "default",
		},
		Spec: spinv1alpha1.SpinAppSpec{
			Image:    "ghcr.io/spinkube/containerd-shim-spin/examples/spin-rust-hello:v0.13.0",
			Replicas: -1,
		},
	})
	require.Error(t, err)
}
