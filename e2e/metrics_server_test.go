package e2e

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"sigs.k8s.io/e2e-framework/support/utils"
)

const metricsServiceURL = "https://spin-operator-metrics-service.spin-operator.svc.cluster.local:8443/metrics"

// TestMetricsServer is a test that checks that the metrics server is
// up and running and protected by authn/authz.
func TestMetricsServer(t *testing.T) {
	var client klient.Client

	defaultTest := features.New("metrics server is functional and protected").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client = cfg.Client()

			if err := client.Resources().Create(ctx, newMetricsClusterRoleBinding()); err != nil {
				t.Fatalf("Failed to create clusterrolebinding: %s", err)
			}

			deployment := newCurlDeployment(testNamespace, "curl-deploy", 1, "curl")
			if err := client.Resources().Create(ctx, deployment); err != nil {
				t.Fatalf("Failed to create curl deployment: %s", err)
			}
			err := wait.For(conditions.New(client.Resources()).DeploymentConditionMatch(deployment, appsv1.DeploymentAvailable, corev1.ConditionTrue), wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatalf("Timed out waiting for curl deployment: %s", err)
			}

			return ctx
		}).
		Assess("metrics service is available",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				svc := &corev1.ServiceList{
					Items: []corev1.Service{
						{ObjectMeta: metav1.ObjectMeta{Name: "spin-operator-metrics-service", Namespace: "spin-operator"}},
					},
				}

				if err := wait.For(
					conditions.New(client.Resources()).ResourcesFound(svc),
					wait.WithTimeout(3*time.Minute),
					wait.WithInterval(time.Second),
				); err != nil {
					t.Fatalf("Timed out waiting for the spin-operator-metrics-service svc resource: %s", err)
				}

				return ctx
			}).
		Assess("metrics server is protected", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			pods := &corev1.PodList{}
			listOpts := func(opts *metav1.ListOptions) { opts.LabelSelector = "app=curl" }
			err := client.Resources(testNamespace).List(ctx, pods, listOpts)
			if err != nil || pods.Items == nil {
				t.Fatalf("error while getting pods: %s", err)
			}
			podName := pods.Items[0].Name

			command := []string{"curl", "-v", "-k", "-H", "Authorization: Bearer bogus", metricsServiceURL}
			var stdout, stderr bytes.Buffer
			if err := client.Resources().ExecInPod(ctx, testNamespace, podName, "curl", command, &stdout, &stderr); err != nil {
				t.Log(stderr.String())
				t.Fatalf("Failed executing the command in the curl pod: %s", err)
			}

			if !strings.Contains(stdout.String(), "Authentication failed") {
				t.Fatalf("Failed to get authentication error with bogus token.\nstdout: %s\nstderr: %s", stdout.String(), stderr.String())
			}

			return ctx
		}).
		Assess("metrics server works when authorized", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			token := utils.FetchCommandOutput("kubectl create token spin-operator-controller-manager -n spin-operator")
			if token == "" {
				t.Fatal("Failed to create token via service account")
			}

			pods := &corev1.PodList{}
			listOpts := func(opts *metav1.ListOptions) { opts.LabelSelector = "app=curl" }
			err := client.Resources(testNamespace).List(ctx, pods, listOpts)
			if err != nil || pods.Items == nil {
				t.Fatalf("error while getting pods: %s", err)
			}
			podName := pods.Items[0].Name

			command := []string{"curl", "-v", "-k", "-H", fmt.Sprintf("Authorization: Bearer %s", token), metricsServiceURL}
			var stdout, stderr bytes.Buffer
			if err := client.Resources().ExecInPod(ctx, testNamespace, podName, "curl", command, &stdout, &stderr); err != nil {
				t.Log(stderr.String())
				t.Fatalf("Failed executing the command in the curl pod: %s", err)
			}

			if !strings.Contains(stdout.String(), "controller_runtime_reconcile_total") {
				t.Fatalf("Failed to find sample metric in curl output.\nstdout: %s\nstderr: %s", stdout.String(), stderr.String())
			}

			return ctx
		}).
		Feature()
	testEnv.Test(t, defaultTest)
}

func newMetricsClusterRoleBinding() *rbac.ClusterRoleBinding {
	crb := rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "spin-operator-metrics-binding",
		},
		Subjects: []rbac.Subject{
			{
				Kind:      rbac.ServiceAccountKind,
				Name:      "spin-operator-controller-manager",
				Namespace: "spin-operator",
			},
		},
		RoleRef: rbac.RoleRef{
			Kind: "ClusterRole",
			Name: "spin-operator-metrics-reader",
		},
	}
	return &crb
}

func newCurlDeployment(namespace string, name string, replicas int32, containerName string) *appsv1.Deployment {
	labels := map[string]string{"app": "curl"}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    containerName,
							Image:   "curlimages/curl:latest",
							Command: []string{"/bin/sleep"},
							Args:    []string{"1000"},
						},
					},
				},
			},
		},
	}
}
