package e2e

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	spinapps_v1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
)

var runtimeClassName = "wasmtime-spin-v2"

// TestDefaultSetup is a test that checks that the minimal setup works
// with the containerd wasm shim runtime as the default runtime.
func TestDefaultSetup(t *testing.T) {
	var client klient.Client

	helloWorldImage := "ghcr.io/spinkube/containerd-shim-spin/examples/spin-rust-hello:v0.13.0"
	testSpinAppName := "test-spinapp"

	defaultTest := features.New("default and most minimal setup").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client = cfg.Client()

			testSpinApp := newSpinAppCR(testSpinAppName, helloWorldImage, "containerd-shim-spin")
			if err := client.Resources().Create(ctx, testSpinApp); err != nil {
				t.Fatalf("Failed to create spinapp: %s", err)
			}

			return ctx
		}).
		Assess("spin app deployment is created and available",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				if err := wait.For(
					conditions.New(client.Resources()).DeploymentAvailable(testSpinAppName, testNamespace),
					wait.WithTimeout(3*time.Minute),
					wait.WithInterval(time.Second),
				); err != nil {
					t.Fatal(err)
				}

				return ctx
			}).
		Assess("spin app service is created and available", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			svc := &corev1.ServiceList{
				Items: []corev1.Service{
					{ObjectMeta: metav1.ObjectMeta{Name: testSpinAppName, Namespace: testNamespace}},
				},
			}

			if err := wait.For(
				conditions.New(client.Resources()).ResourcesFound(svc),
				wait.WithTimeout(3*time.Minute),
				wait.WithInterval(time.Second),
			); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).
		Feature()
	testEnv.Test(t, defaultTest)
}

func newSpinAppCR(name, image, executor string) *spinapps_v1alpha1.SpinApp {
	return &spinapps_v1alpha1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Spec: spinapps_v1alpha1.SpinAppSpec{
			Replicas: 1,
			Image:    image,
			Executor: executor,
		},
	}
}
