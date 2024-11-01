// helper is a package that offers e2e test helpers. It's a badly named package.
// If it grows we should refactor it into something a little more manageable
// but it's fine for now.
package helper

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	controllerruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

const debugDeploymentName = "debugy"

// EnsureDebugContainer ensures that the helper debug container is installed right namespace. This allows us to make requests from inside the cluster
// regardless of the external network configuration.
func EnsureDebugContainer(t *testing.T, ctx context.Context, cfg *envconf.Config, namespace string) {
	t.Helper()

	client, err := cfg.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	// Deploy a debug container so that we can test that the app is available later
	deployment := newDebugDeployment(namespace, debugDeploymentName, 1, debugDeploymentName)
	if err = client.Resources().Create(ctx, deployment); controllerruntimeclient.IgnoreAlreadyExists(err) != nil {
		t.Fatal(err)
	}

	err = wait.For(
		conditions.New(client.Resources()).
			DeploymentConditionMatch(deployment, appsv1.DeploymentAvailable, corev1.ConditionTrue),
		wait.WithTimeout(time.Minute*5))
	if err != nil {
		t.Fatal(err)
	}
}

// CurlSpinApp is a crude function for using the debug pod to send a HTTP request to a spin app
// within the cluster. It allows customization of the route, and all requests are GET requests unless a non-empty body is provided.
func CurlSpinApp(t *testing.T, ctx context.Context, cfg *envconf.Config, namespace, spinAppName, route, body string) (string, int, error) {
	t.Helper()

	client, err := cfg.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	// Find the debug pod
	pods := &corev1.PodList{}
	err = client.Resources(namespace).List(ctx, pods, resources.WithLabelSelector("app=pod-exec"))
	if err != nil || pods.Items == nil || len(pods.Items) == 0 {
		return "", -1, fmt.Errorf("failed to get debug pods: %w", err)
	}

	debugPod := pods.Items[0]

	podName := debugPod.Name

	command := []string{"curl", "--silent", "--max-time", "5", "--write-out", "\n%{http_code}\n", "http://" + spinAppName + "." + namespace + route, "--output", "-"}
	if body != "" {
		command = append(command, "--data", body)
	}

	var stdout, stderr bytes.Buffer
	if err := client.Resources().ExecInPod(ctx, namespace, podName, debugDeploymentName, command, &stdout, &stderr); err != nil {
		t.Logf("Curl Spin App failed, err: %v.\nstdout:\n%s\n\nstderr:\n%s\n", err, stdout.String(), stderr.String())
		return "", -1, err
	}

	parts := strings.SplitN(stdout.String(), "\n", 2)
	if len(parts) != 2 {
		t.Fatalf("Curl Spin App failed, unexpected response format: %s", &stdout)
	}

	strStatus := strings.Trim(parts[1], "\n")
	statusCode, err := strconv.Atoi(strStatus)
	if err != nil {
		t.Logf("error parsing status code: %v", err)
		return parts[0], statusCode, err
	}
	t.Logf("Curl Spin App response: %s, status code: %d, err: %v", parts[0], statusCode, err)
	return parts[0], statusCode, nil

}

func newDebugDeployment(namespace string, name string, replicas int32, containerName string) *appsv1.Deployment {
	labels := map[string]string{"app": "pod-exec"}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: containerName, Image: "nginx"}}},
			},
		},
	}
}
