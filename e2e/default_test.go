package e2e

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	nodev1 "k8s.io/api/node/v1"
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

var runtimeClassName = "wasmtime-spin-v2"

const (
	testPodLabelName  = "core.spinoperator.dev/test"
	testPodLabelValue = "foobar"
)

// TestDefaultSetup is a test that checks that the minimal setup works
// with the containerd wasm shim runtime as the default runtime.
func TestDefaultSetup(t *testing.T) {
	var client klient.Client

	helloWorldImage := "ghcr.io/spinkube/containerd-shim-spin/examples/spin-rust-hello:v0.13.0"
	testSpinAppName := "test-spinapp"

	defaultTest := features.New("default and most minimal setup").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

			client = cfg.Client()

			if err := spinapps_v1alpha1.AddToScheme(client.Resources(testNamespace).GetScheme()); err != nil {
				t.Fatalf("failed to register the spinapps_v1alpha1 types with Kuberenets scheme: %s", err)
			}

			runtimeClass := &nodev1.RuntimeClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: runtimeClassName,
				},
				Handler: "spin",
			}

			if err := client.Resources().Create(ctx, runtimeClass); err != nil {
				t.Fatalf("Failed to create runtimeclass: %s", err)
			}

			return ctx
		}).
		Assess("spin app custom resource is created", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			testSpinApp := newSpinAppCR(testSpinAppName, helloWorldImage)

			if err := client.Resources().Create(ctx, newContainerdShimExecutor(testNamespace)); err != nil {
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

			if err := wait.For(
				conditions.New(client.Resources()).ResourcesFound(svc),
				wait.WithTimeout(3*time.Minute),
				wait.WithInterval(30*time.Second),
			); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()
	testEnv.Test(t, defaultTest)
}

func newSpinAppCR(name, image string) *spinapps_v1alpha1.SpinApp {
	var testSpinApp = &spinapps_v1alpha1.SpinApp{
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

	return testSpinApp

}

func newContainerdShimExecutor(namespace string) *spinapps_v1alpha1.SpinAppExecutor {
	var testSpinAppExecutor = &spinapps_v1alpha1.SpinAppExecutor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerd-shim-spin",
			Namespace: namespace,
		},
		Spec: spinapps_v1alpha1.SpinAppExecutorSpec{
			CreateDeployment: true,
			DeploymentConfig: &spinapps_v1alpha1.ExecutorDeploymentConfig{
				RuntimeClassName: runtimeClassName,
			},
		},
	}

	return testSpinAppExecutor
}
