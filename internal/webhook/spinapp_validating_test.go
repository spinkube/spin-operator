package webhook

import (
	"errors"
	"testing"

	spinv1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/internal/constants"
	"github.com/stretchr/testify/require"
)

func TestValidateExecutor(t *testing.T) {
	t.Parallel()

	_, fldErr := validateExecutor(spinv1alpha1.SpinAppSpec{}, func(string) (*spinv1alpha1.SpinAppExecutor, error) { return nil, nil })
	require.EqualError(t, fldErr, "spec.executor: Invalid value: \"\": executor must be set, likely no default executor was set because you have no executors installed")

	_, fldErr = validateExecutor(
		spinv1alpha1.SpinAppSpec{Executor: constants.CyclotronExecutor},
		func(string) (*spinv1alpha1.SpinAppExecutor, error) { return nil, errors.New("executor not found?") })
	require.EqualError(t, fldErr, "spec.executor: Invalid value: \"cyclotron\": executor does not exist on cluster")

	_, fldErr = validateExecutor(spinv1alpha1.SpinAppSpec{Executor: constants.ContainerDShimSpinExecutor}, func(string) (*spinv1alpha1.SpinAppExecutor, error) { return nil, nil })
	require.Nil(t, fldErr)
}

func TestValidateAnnotations(t *testing.T) {
	t.Parallel()

	deploymentlessExecutor := &spinv1alpha1.SpinAppExecutor{
		Spec: spinv1alpha1.SpinAppExecutorSpec{
			CreateDeployment: false,
		},
	}
	deploymentfullExecutor := &spinv1alpha1.SpinAppExecutor{
		Spec: spinv1alpha1.SpinAppExecutorSpec{
			CreateDeployment: true,
		},
	}

	fldErr := validateAnnotations(spinv1alpha1.SpinAppSpec{
		Executor:              "an-executor",
		DeploymentAnnotations: map[string]string{"key": "asdf"},
	}, deploymentlessExecutor)
	require.EqualError(t, fldErr,
		`spec.deploymentAnnotations: Invalid value: map[string]string{"key":"asdf"}: `+
			`deploymentAnnotations can't be set when the executor does not use operator deployments`)

	fldErr = validateAnnotations(spinv1alpha1.SpinAppSpec{
		Executor:       "an-executor",
		PodAnnotations: map[string]string{"key": "asdf"},
	}, deploymentlessExecutor)
	require.EqualError(t, fldErr,
		`spec.podAnnotations: Invalid value: map[string]string{"key":"asdf"}: `+
			`podAnnotations can't be set when the executor does not use operator deployments`)

	fldErr = validateAnnotations(spinv1alpha1.SpinAppSpec{
		Executor:              "an-executor",
		DeploymentAnnotations: map[string]string{"key": "asdf"},
	}, deploymentfullExecutor)
	require.Nil(t, fldErr)

	fldErr = validateAnnotations(spinv1alpha1.SpinAppSpec{
		Executor: "an-executor",
	}, deploymentlessExecutor)
	require.Nil(t, fldErr)
}
