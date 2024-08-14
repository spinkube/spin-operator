package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	spinv1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/pkg/spinapp"
)

func minimalSpinApp() *spinv1alpha1.SpinApp {
	return &spinv1alpha1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-app",
			Namespace: "default",
		},
		Spec: spinv1alpha1.SpinAppSpec{
			Executor: "containerd-shim-spin",
			Image:    "fakereg.dev/noapp:latest",
			Replicas: 1,
		},
	}
}

func TestConstructRuntimeConfigSecretMount_Contract(t *testing.T) {
	t.Parallel()

	volume, mount := constructRuntimeConfigSecretMount(context.Background(), "my-secret-v1")
	// We currently expect runtime config to be optional.
	// TODO: evaluate whether we should require this - silently not loading config
	//       feels subpar.
	require.True(t, *volume.VolumeSource.Secret.Optional)

	// Require the volume to be spin- prefixed to avoid collisions with user volumes.
	require.Contains(t, volume.Name, "spin-")

	// Require the volume mount to be spin- prefixed to avoid collisions with user volumes.
	require.Contains(t, mount.Name, "spin-")
}

func TestConstructCASecretMount(t *testing.T) {
	t.Parallel()

	volume, mount := constructCASecretMount(context.Background(), "a-secret-name")

	// Mount and Volume refer to each other
	require.Equal(t, volume.Name, mount.Name)

	// uses provided secret name
	require.Equal(t, "a-secret-name", volume.Secret.SecretName)

	require.True(t, mount.ReadOnly)
}

func TestConstructVolumeMountsForApp_Contract(t *testing.T) {
	t.Parallel()

	// There should be an error when trying to load runtime-config from multiple
	// places.
	app := minimalSpinApp()
	app.Spec.RuntimeConfig.LoadFromSecret = "a-secret"
	_, _, err := ConstructVolumeMountsForApp(context.Background(), app, "a-generated-secret", "a-ca-secret")
	require.Error(t, err)
	require.ErrorContains(t, err, "cannot specify both a user-provided runtime secret and a generated one")

	// No runtime secret at all is ok
	app = minimalSpinApp()
	app.Spec.RuntimeConfig.LoadFromSecret = ""
	volumes, mounts, err := ConstructVolumeMountsForApp(context.Background(), app, "", "")
	require.NoError(t, err)
	require.Len(t, volumes, 0)
	require.Len(t, mounts, 0)

	// User provided runtime secret is ok
	app = minimalSpinApp()
	app.Spec.RuntimeConfig.LoadFromSecret = "foo-secret-v1"
	volumes, mounts, err = ConstructVolumeMountsForApp(context.Background(), app, "", "")
	require.NoError(t, err)
	require.Len(t, volumes, 1)
	require.Len(t, mounts, 1)
	require.Equal(t, "foo-secret-v1", volumes[0].VolumeSource.Secret.SecretName)

	// Generated runtime secret is ok
	app = minimalSpinApp()
	app.Spec.RuntimeConfig.LoadFromSecret = ""
	volumes, mounts, err = ConstructVolumeMountsForApp(context.Background(), app, "gen-secret", "spin-ca")
	require.NoError(t, err)
	require.Len(t, volumes, 2)
	require.Len(t, mounts, 2)
	require.Equal(t, "gen-secret", volumes[0].VolumeSource.Secret.SecretName)
}

func TestConstructEnvForApp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		varName         string
		expectedEnvName string

		value     string
		valueFrom *corev1.EnvVarSource
	}{
		{
			name:            "simple_secret_with_static_value",
			varName:         "simple_secret",
			expectedEnvName: "SPIN_VARIABLE_SIMPLE_SECRET",
			value:           "f00",
		},
		{
			name:            "simple_secret_with_numb3rs_and_static_value",
			varName:         "simple_secret_with_numb3rs",
			expectedEnvName: "SPIN_VARIABLE_SIMPLE_SECRET_WITH_NUMB3RS",
			value:           "f00",
		},
		{
			name:            "simple_secret_with_secret_value",
			varName:         "simple_secret",
			expectedEnvName: "SPIN_VARIABLE_SIMPLE_SECRET",
			valueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "my-secret",
					},
				},
			},
		},
		{
			name:            "pod_attribute_value",
			varName:         "pod_namespace",
			expectedEnvName: "SPIN_VARIABLE_POD_NAMESPACE",
			valueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := minimalSpinApp()
			app.Spec.Variables = []spinv1alpha1.SpinVar{
				{
					Name:      test.varName,
					Value:     test.value,
					ValueFrom: test.valueFrom,
				},
			}

			envs := ConstructEnvForApp(context.Background(), app, 0)

			require.Equal(t, test.expectedEnvName, envs[0].Name)
			require.Equal(t, test.value, envs[0].Value)
			require.Equal(t, test.valueFrom, envs[0].ValueFrom)
		})
	}
}

func TestSpinHealthCheckToCoreProbe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		probe         *spinv1alpha1.HealthProbe
		expectedProbe *corev1.Probe
		expectedErr   string
	}{
		{
			name:          "no_probe",
			probe:         nil,
			expectedProbe: nil,
		},
		{
			name:          "probe_missing_httpGet_spec",
			probe:         &spinv1alpha1.HealthProbe{},
			expectedProbe: nil,
			expectedErr:   "probe exists but with unknown configuration",
		},
		{
			name: "probe_full",
			probe: &spinv1alpha1.HealthProbe{
				HTTPGet: &spinv1alpha1.HTTPHealthProbe{
					Path: "/var",
					HTTPHeaders: []spinv1alpha1.HTTPHealthProbeHeader{
						{
							Name:  "header",
							Value: "value",
						},
					},
				},
				InitialDelaySeconds: 1,
				TimeoutSeconds:      2,
				PeriodSeconds:       3,
				SuccessThreshold:    4,
				FailureThreshold:    5,
			},
			expectedProbe: &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/var",
						Port: intstr.FromInt(80),
						HTTPHeaders: []corev1.HTTPHeader{
							{
								Name:  "header",
								Value: "value",
							},
						},
					},
				},
				InitialDelaySeconds: 1,
				TimeoutSeconds:      2,
				PeriodSeconds:       3,
				SuccessThreshold:    4,
				FailureThreshold:    5,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := SpinHealthCheckToCoreProbe(test.probe)
			if test.expectedErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, test.expectedErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.expectedProbe, result)
		})
	}
}

func TestDeploymentLabel(t *testing.T) {
	scheme := registerAndGetScheme()
	app := minimalSpinApp()
	deployment, err := constructDeployment(context.Background(), app, &spinv1alpha1.ExecutorDeploymentConfig{}, "", "", scheme)

	require.Nil(t, err)
	require.NotNil(t, deployment.ObjectMeta.Labels)
	require.Equal(t, deployment.ObjectMeta.Labels[spinapp.NameLabelKey], app.Name)
}

func registerAndGetScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(spinv1alpha1.AddToScheme(scheme))

	return scheme
}

func spinAppWithLabels(labels map[string]string) *spinv1alpha1.SpinApp {
	return &spinv1alpha1.SpinApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-app",
			Namespace: "default",
		},
		Spec: spinv1alpha1.SpinAppSpec{
			Executor:  "containerd-shim-spin",
			Image:     "fakereg.dev/noapp:latest",
			Replicas:  1,
			PodLabels: labels,
		},
	}
}
