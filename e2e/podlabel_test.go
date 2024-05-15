package e2e

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	spinapps_v1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
)

const (
	testPodLabelName  = "core.spinoperator.dev/test"
	testPodLabelValue = "foobar"
)

// TestPodLabels is a test that checks that SpinApp.spec.podLabels are
// passed down to underlying pods
func TestPodLabels(t *testing.T) {
	var client klient.Client

	helloWorldImage := "ghcr.io/spinkube/containerd-shim-spin/examples/spin-rust-hello:v0.13.0"
	testSpinAppName := "test-spinapp-with-pod-labels"

	podLabelTest := features.New("SpinApp with custom PodLabels").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

			client = cfg.Client()

			if err := spinapps_v1alpha1.AddToScheme(client.Resources(testNamespace).GetScheme()); err != nil {
				t.Fatalf("failed to register the spinapps_v1alpha1 types with Kuberenets scheme: %s", err)
			}

			return ctx
		}).
		Assess("spin app custom resource is created", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			testSpinApp := newSpinAppWithPodLabels(testSpinAppName, helloWorldImage)

			if err := client.Resources().Create(ctx, testSpinApp); err != nil {
				t.Fatalf("Failed to create spinapp: %s", err)
			}

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
		Assess("spin app deployment is available", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			if err := wait.For(
				conditions.New(client.Resources()).DeploymentAvailable(testSpinAppName, testNamespace),
				wait.WithTimeout(3*time.Minute),
				wait.WithInterval(30*time.Second),
			); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).
		Assess("spin app deployment has custom pod labels", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			res, err := resources.New(client.RESTConfig())
			if err != nil {
				t.Fatalf("Could not create controller runtime client: %s", err)
			}
			var deploy appsv1.Deployment
			if err = res.Get(ctx, testSpinAppName, testNamespace, &deploy); err != nil {
				t.Fatalf("Could not find deployment: %s", err)
			}
			v, ok := deploy.Spec.Template.Labels[testPodLabelName]
			if !ok || v != testPodLabelValue {
				t.Fatal("PodLabels were not passed from the SpinApp to the underlying Pod")
			}
			return ctx
		}).
		Feature()
	testEnv.Test(t, podLabelTest)
}

func newSpinAppWithPodLabels(name, image string) *spinapps_v1alpha1.SpinApp {
	return &spinapps_v1alpha1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Spec: spinapps_v1alpha1.SpinAppSpec{
			Replicas:  1,
			Image:     image,
			Executor:  "containerd-shim-spin",
			PodLabels: map[string]string{testPodLabelName: testPodLabelValue},
		},
	}
}
