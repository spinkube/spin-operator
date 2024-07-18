package e2e

import (
	"context"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	spinapps_v1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/internal/generics"
)

// TestSpintainer is a test that checks that the minimal setup works
// with the spintainer executor
func TestSpintainer(t *testing.T) {
	var client klient.Client

	helloWorldImage := "ghcr.io/spinkube/spin-operator/hello-world:20240708-130250-gfefd2b1"
	testSpinAppName := "test-spintainer-app"

	defaultTest := features.New("default and most minimal setup").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

			client = cfg.Client()

			if err := spinapps_v1alpha1.AddToScheme(client.Resources(testNamespace).GetScheme()); err != nil {
				t.Fatalf("failed to register the spinapps_v1alpha1 types with Kubernetes scheme: %s", err)
			}

			return ctx
		}).
		Assess("spin app custom resource is created", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			testSpinApp := newSpinAppCR(testSpinAppName, helloWorldImage, "spintainer")

			if err := client.Resources().Create(ctx, newSpintainerExecutor(testNamespace)); err != nil {
				t.Fatalf("Failed to create spinappexecutor: %s", err)
			}

			if err := client.Resources().Create(ctx, testSpinApp); err != nil {
				t.Fatalf("Failed to create spinapp: %s", err)
			}
			// wait for spinapp to be created
			if err := wait.For(
				conditions.New(client.Resources()).ResourceMatch(testSpinApp, func(object k8s.Object) bool {
					return true
				}),
				wait.WithTimeout(3*time.Minute),
				wait.WithInterval(30*time.Second),
			); err != nil {
				t.Fatal(err)
			}

			return ctx
		}).
		Assess("spin app deployment and service are available", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// wait for deployment to be ready
			if err := wait.For(
				conditions.New(client.Resources()).DeploymentAvailable(testSpinAppName, testNamespace),
				wait.WithTimeout(3*time.Minute),
				wait.WithInterval(30*time.Second),
			); err != nil {
				t.Fatal(err)
			}

			svc := &v1.ServiceList{
				Items: []v1.Service{
					{ObjectMeta: metav1.ObjectMeta{Name: testSpinAppName, Namespace: testNamespace}},
				},
			}

			if err := wait.For(
				conditions.New(client.Resources()).ResourcesFound(svc),
				wait.WithTimeout(3*time.Minute),
				wait.WithInterval(30*time.Second),
			); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).
		Feature()
	testEnv.Test(t, defaultTest)
}

func newSpintainerExecutor(namespace string) *spinapps_v1alpha1.SpinAppExecutor {
	var testSpinAppExecutor = &spinapps_v1alpha1.SpinAppExecutor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spintainer",
			Namespace: namespace,
		},
		Spec: spinapps_v1alpha1.SpinAppExecutorSpec{
			CreateDeployment: true,
			DeploymentConfig: &spinapps_v1alpha1.ExecutorDeploymentConfig{
				SpinImage: generics.Ptr("ghcr.io/fermyon/spin:v2.7.0"),
			},
		},
	}

	return testSpinAppExecutor
}
