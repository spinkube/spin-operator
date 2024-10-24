package e2e

import (
	"context"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	spinapps_v1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/e2e/helper"
	"github.com/stretchr/testify/require"
)

// TestComponentFiltering checks that the operator and shim have support for running a subset of a Spin app's components
func TestComponentFiltering(t *testing.T) {
	var client klient.Client

	// TODO: Use an image from a sample app in this repository
	appImage := "ghcr.io/kate-goldenring/spin-operator/examples/spin-salutations:20241022-144454"
	testSpinAppName := "test-component-filtering"

	defaultTest := features.New("default and most minimal setup").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

			client = cfg.Client()

			if err := spinapps_v1alpha1.AddToScheme(client.Resources(testNamespace).GetScheme()); err != nil {
				t.Fatalf("failed to register the spinapps_v1alpha1 types with Kubernetes scheme: %s", err)
			}

			return ctx
		}).
		Assess("spin app custom resource is created", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			testSpinApp := newSpinAppWithComponentFiltering(testSpinAppName, appImage, "containerd-shim-spin")

			if err := client.Resources().Create(ctx, testSpinApp); err != nil {
				t.Fatalf("Failed to create spinapp: %s", err)
			}
			return ctx
		}).
		Assess("spin app deployment and service are available", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// wait for deployment to be ready
			if err := wait.For(
				conditions.New(client.Resources()).DeploymentAvailable(testSpinAppName, testNamespace),
				wait.WithTimeout(3*time.Minute),
				wait.WithInterval(time.Second),
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
				wait.WithInterval(500*time.Millisecond),
			); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).
		Assess("spin app is only serving hello component", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			helper.EnsureDebugContainer(t, ctx, cfg, testNamespace)

			_, status, err := helper.CurlSpinApp(t, ctx, cfg, testNamespace, testSpinAppName, "/hi", "")

			require.NoError(t, err)
			require.Equal(t, 200, status)

			_, status, err = helper.CurlSpinApp(t, ctx, cfg, testNamespace, testSpinAppName, "/bye", "")

			require.NoError(t, err)
			require.Equal(t, 404, status)

			return ctx
		}).
		Feature()
	testEnv.Test(t, defaultTest)
}

// TODO - Create Spintainer once `--component-id` flag is available in Spin 3.0 and support added in Spintainer to pull flags to set from env var
// func TestSpintainerTestComponentFiltering(t *testing.T) {
// }

func newSpinAppWithComponentFiltering(name, image, executor string) *spinapps_v1alpha1.SpinApp {
	return &spinapps_v1alpha1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Spec: spinapps_v1alpha1.SpinAppSpec{
			Replicas: 1,
			Image:    image,
			Executor: executor,
			// Only execute the hello component
			Components: []string{"hello"},
		},
	}
}
