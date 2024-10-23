package e2e

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	controllerruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"sigs.k8s.io/e2e-framework/support/utils"

	spinapps_v1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/e2e/helper"
	"github.com/spinkube/spin-operator/internal/generics"
	"github.com/stretchr/testify/require"
)

// TestShimRedis checks that the shim has basic runtime config support
// by being able to integrate with redis as the backend for a kv store
func TestShimRedis(t *testing.T) {
	var client klient.Client

	// TODO: Use an image from a sample app in this repository
	appImage := "ghcr.io/calebschoepp/spin-checklist:v0.1.0"
	testSpinAppName := "test-shim-redis"

	defaultTest := features.New("default and most minimal setup").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

			client = cfg.Client()

			if err := spinapps_v1alpha1.AddToScheme(client.Resources(testNamespace).GetScheme()); err != nil {
				t.Fatalf("failed to register the spinapps_v1alpha1 types with Kubernetes scheme: %s", err)
			}

			return ctx
		}).
		Assess("redis is setup", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			if err := client.Resources().Create(ctx, redisDeployment()); err != nil {
				t.Fatalf("Failed to create redis deployment: %s", err)
			}

			if err := client.Resources().Create(ctx, redisService()); err != nil {
				t.Fatalf("Failed to create redis service: %s", err)
			}

			return ctx
		}).
		Assess("spin app custom resource is created", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			testSpinApp := newSpinAppUsingRedis(testSpinAppName, appImage, "containerd-shim-spin")

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
		Assess("spin app is using redis", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			helper.EnsureDebugContainer(t, ctx, cfg, testNamespace)

			_, status, err := helper.CurlSpinApp(t, ctx, cfg, testNamespace, testSpinAppName, "/api", "{\"key\": \"foo\", \"value\": \"bar\"}")

			require.NoError(t, err)
			require.Equal(t, 200, status)

			p := utils.RunCommand("kubectl exec -n default deployment/redis -- bash -c \"echo KEYS * | redis-cli\"")
			require.NoError(t, p.Err())
			buf := new(strings.Builder)
			_, err = io.Copy(buf, p.Out())
			require.NoError(t, err)
			if !strings.Contains(buf.String(), "foo") {
				t.Fatalf("expected 'foo' but got '%s'", buf.String())
			}

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			if err := client.Resources().Delete(ctx, redisDeployment()); err != nil {
				t.Fatalf("Failed to delete redis deployment: %s", err)
			}
			if err := client.Resources().Delete(ctx, redisService()); err != nil {
				t.Fatalf("Failed to delete redis service: %s", err)
			}
			return ctx
		}).
		Feature()
	testEnv.Test(t, defaultTest)
}

// TestSpintainerRedis checks that Spintainer has basic runtime config support
// by being able to integrate with redis as the backend for a kv store
//
//nolint:gocyclo
func TestSpintainerRedis(t *testing.T) {
	var client klient.Client

	// TODO: Use an image from a sample app in this repository
	appImage := "ghcr.io/calebschoepp/spin-checklist:v0.1.0"
	testSpinAppName := "test-spintainer-redis"

	defaultTest := features.New("default and most minimal setup").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

			client = cfg.Client()

			if err := spinapps_v1alpha1.AddToScheme(client.Resources(testNamespace).GetScheme()); err != nil {
				t.Fatalf("failed to register the spinapps_v1alpha1 types with Kubernetes scheme: %s", err)
			}

			return ctx
		}).
		Assess("redis is setup", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			if err := client.Resources().Create(ctx, redisDeployment()); err != nil {
				t.Fatalf("Failed to create redis deployment: %s", err)
			}

			if err := client.Resources().Create(ctx, redisService()); err != nil {
				t.Fatalf("Failed to create redis service: %s", err)
			}

			return ctx
		}).
		Assess("spin app custom resource is created", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			testSpinApp := newSpinAppUsingRedis(testSpinAppName, appImage, "spintainer")

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
		Assess("spin app is using redis", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			helper.EnsureDebugContainer(t, ctx, cfg, testNamespace)

			_, status, err := helper.CurlSpinApp(t, ctx, cfg, testNamespace, testSpinAppName, "/api", "{\"key\": \"foo\", \"value\": \"bar\"}")

			require.NoError(t, err)
			require.Equal(t, 200, status)

			p := utils.RunCommand("kubectl exec -n default deployment/redis -- bash -c \"echo KEYS * | redis-cli\"")
			require.NoError(t, p.Err())
			buf := new(strings.Builder)
			_, err = io.Copy(buf, p.Out())
			require.NoError(t, err)
			if !strings.Contains(buf.String(), "foo") {
				t.Fatalf("expected 'foo' but got '%s'", buf.String())
			}

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			if err := client.Resources().Delete(ctx, redisDeployment()); err != nil {
				t.Fatalf("Failed to delete redis deployment: %s", err)
			}
			if err := client.Resources().Delete(ctx, redisService()); err != nil {
				t.Fatalf("Failed to delete redis service: %s", err)
			}
			return ctx
		}).
		Feature()
	testEnv.Test(t, defaultTest)
}

func redisDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: generics.Ptr(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "redis"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "redis"},
				},
				Spec: v1.PodSpec{
					Volumes:        []v1.Volume{},
					InitContainers: []v1.Container{},
					Containers: []v1.Container{{
						Name:  "redis",
						Image: "redis:latest",
						Ports: []v1.ContainerPort{{
							ContainerPort: 6379,
						}},
					}},
				},
			},
		},
	}
}

func redisService() *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"app": "redis"},
			Ports: []v1.ServicePort{{
				Protocol:   v1.ProtocolTCP,
				Port:       6379,
				TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 6379},
			}},
		},
	}
}

func newSpinAppUsingRedis(name, image, executor string) *spinapps_v1alpha1.SpinApp {
	return &spinapps_v1alpha1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Spec: spinapps_v1alpha1.SpinAppSpec{
			Replicas: 1,
			Image:    image,
			Executor: executor,
			RuntimeConfig: spinapps_v1alpha1.RuntimeConfig{
				KeyValueStores: []spinapps_v1alpha1.KeyValueStoreConfig{{
					Name: "default",
					Type: "redis",
					Options: []spinapps_v1alpha1.RuntimeConfigOption{{
						Name:  "url",
						Value: "redis://redis.default.svc.cluster.local:6379",
					}},
				}},
			},
		},
	}
}
