package e2e

import (
	"context"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	controllerruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	spinapps_v1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/e2e/helper"
	"github.com/stretchr/testify/require"
)

// TestSelectiveDeployment checks that the operator and shim have support for running a subset of a Spin app's components
func TestSelectiveDeployment(t *testing.T) {
	var client klient.Client

	appImage := "ghcr.io/spinkube/spin-operator/salutations:20241105-223428-g4da3171"
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
			testSpinApp := newSpinAppCR(testSpinAppName, appImage, "containerd-shim-spin", []string{"hello"})

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
			return testSelectiveDeploymentOfSalutationsApp(ctx, t, cfg, testSpinAppName)
		}).
		Feature()
	testEnv.Test(t, defaultTest)
}

// TestSelectiveDeploymentSpintainer is a test that checks that selective deployment works
func TestSelectiveDeploymentSpintainer(t *testing.T) {
	var client klient.Client

	salutationsApp := "ghcr.io/spinkube/spin-operator/salutations:20241105-223428-g4da3171"
	testSpinAppName := "test-spintainer-selective-deployment"

	defaultTest := features.New("default and most minimal setup").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

			client = cfg.Client()

			if err := spinapps_v1alpha1.AddToScheme(client.Resources(testNamespace).GetScheme()); err != nil {
				t.Fatalf("failed to register the spinapps_v1alpha1 types with Kubernetes scheme: %s", err)
			}

			return ctx
		}).
		Assess("spin app custom resource is created", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			testSpinApp := newSpinAppCR(testSpinAppName, salutationsApp, "spintainer", []string{"hello"})

			if err := client.Resources().Create(ctx, newSpintainerExecutor(testNamespace)); controllerruntimeclient.IgnoreAlreadyExists(err) != nil {
				t.Fatalf("Failed to create spinappexecutor: %s", err)
			}

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
			return testSelectiveDeploymentOfSalutationsApp(ctx, t, cfg, testSpinAppName)
		}).
		Feature()
	testEnv.Test(t, defaultTest)
}

func testSelectiveDeploymentOfSalutationsApp(ctx context.Context, t *testing.T, cfg *envconf.Config, testSpinAppName string) context.Context {
	helper.EnsureDebugContainer(t, ctx, cfg, testNamespace)

	_, status, err := helper.CurlSpinApp(t, ctx, cfg, testNamespace, testSpinAppName, "/hi", "")

	require.NoError(t, err)
	require.Equal(t, 200, status)

	_, status, err = helper.CurlSpinApp(t, ctx, cfg, testNamespace, testSpinAppName, "/bye", "")

	require.NoError(t, err)
	require.Equal(t, 404, status)

	return ctx
}
