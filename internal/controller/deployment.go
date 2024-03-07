package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	spinv1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/internal/generics"
	"github.com/spinkube/spin-operator/pkg/spinapp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func constructRuntimeConfigSecretMount(_ctx context.Context, secretName string) (corev1.Volume, corev1.VolumeMount) {
	volume := corev1.Volume{
		Name: "spin-runtime-config",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
				Optional:   ptr(true),
				Items: []corev1.KeyToPath{
					{
						Key:  "runtime-config.toml",
						Path: "runtime-config.toml",
					},
				},
			},
		},
	}
	volumeMount := corev1.VolumeMount{
		Name:      "spin-runtime-config",
		ReadOnly:  true,
		MountPath: "/runtime-config.toml",
		SubPath:   "runtime-config.toml",
	}

	return volume, volumeMount
}

// ConstructVolumeMountsForApp introspects the application and generates
// any required volume mounts. A generated runtime secret is mutually
// exclusive with a user-provided secret - this is to require _either_ a
// manual runtime-config or a generated one from the crd.
func ConstructVolumeMountsForApp(ctx context.Context, app *spinv1alpha1.SpinApp, generatedRuntimeSecret string) ([]corev1.Volume, []corev1.VolumeMount, error) {
	volumes := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}

	userProvidedRuntimeSecret := app.Spec.RuntimeConfig.LoadFromSecret
	if userProvidedRuntimeSecret != "" && generatedRuntimeSecret != "" {
		return nil, nil, errors.New("cannot specify both a user-provided runtime secret and a generated one")
	}

	selectedSecret := userProvidedRuntimeSecret
	if generatedRuntimeSecret != "" {
		selectedSecret = generatedRuntimeSecret
	}

	if selectedSecret != "" {
		runtimeConfigVolume, runtimeConfigMount := constructRuntimeConfigSecretMount(ctx, selectedSecret)
		volumes = append(volumes, runtimeConfigVolume)
		volumeMounts = append(volumeMounts, runtimeConfigMount)
	}

	// TODO: Once #49 lands validate that volumes don't start with `spin-` prefix in admission webhook.
	volumes = append(volumes, app.Spec.Volumes...)
	volumeMounts = append(volumeMounts, app.Spec.VolumeMounts...)

	return volumes, volumeMounts, nil
}

// ConstructEnvForApp constructs the env for a spin app that runs as a k8s pod.
// Variables are not guaranteed to stay backed by ENV.
func ConstructEnvForApp(ctx context.Context, app *spinv1alpha1.SpinApp) []corev1.EnvVar {
	if len(app.Spec.Variables) == 0 {
		return nil
	}

	envs := make([]corev1.EnvVar, len(app.Spec.Variables))
	for idx, variable := range app.Spec.Variables {
		env := corev1.EnvVar{
			// Spin Variables only allow lowercase ascii characters, `_`, and numbers.
			// this means that we can do a relatively simple conversion here and in
			// the future should implement stronger validation in the webhook/crd definition.
			Name:      fmt.Sprintf("SPIN_VARIABLE_%s", strings.ToUpper(variable.Name)),
			Value:     variable.Value,
			ValueFrom: variable.ValueFrom,
		}
		envs[idx] = env
	}

	return envs
}

func SpinHealthCheckToCoreProbe(probe *spinv1alpha1.HealthProbe) (*corev1.Probe, error) {
	if probe == nil {
		return nil, nil
	}

	if probe.HTTPGet == nil {
		// When the probe is specified, but httpGet is nil, we probably updated the CRD
		// without updating the code. This error is a little janky, but shouldn't ever be seen by
		// an end user.
		return nil, errors.New("probe exists but with unknown configuration, expected httpGet")
	}

	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: probe.HTTPGet.Path,
				Port: intstr.FromInt(spinapp.DefaultHTTPPort),
				HTTPHeaders: generics.MapList(probe.HTTPGet.HTTPHeaders, func(h spinv1alpha1.HTTPHealthProbeHeader) corev1.HTTPHeader {
					return corev1.HTTPHeader{
						Name:  h.Name,
						Value: h.Value,
					}
				}),
			},
		},
		InitialDelaySeconds: probe.InitialDelaySeconds,
		TimeoutSeconds:      probe.TimeoutSeconds,
		PeriodSeconds:       probe.PeriodSeconds,
		SuccessThreshold:    probe.SuccessThreshold,
		FailureThreshold:    probe.FailureThreshold,
	}, nil
}

func ConstructPodHealthChecks(app *spinv1alpha1.SpinApp) (readiness *corev1.Probe, liveness *corev1.Probe, err error) {
	if app.Spec.Checks.Readiness == nil && app.Spec.Checks.Liveness == nil {
		return nil, nil, nil
	}

	readiness, err = SpinHealthCheckToCoreProbe(app.Spec.Checks.Readiness)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to construct readiness probe: %w", err)
	}

	liveness, err = SpinHealthCheckToCoreProbe(app.Spec.Checks.Liveness)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to construct liveness probe: %w", err)
	}

	return readiness, liveness, nil
}
