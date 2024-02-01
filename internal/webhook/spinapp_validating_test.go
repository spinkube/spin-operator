package webhook

import (
	"testing"

	spinv1 "github.com/spinkube/spin-operator/api/v1"
	"github.com/spinkube/spin-operator/internal/constants"
	"github.com/stretchr/testify/require"
)

func TestValidateExecutor(t *testing.T) {
	t.Parallel()

	fldErr := validateExecutor(spinv1.SpinAppSpec{}, func(string) bool { return true })
	require.EqualError(t, fldErr, "spec.executor: Invalid value: \"\": executor must be set, likely no default executor was set because you have no executors installed")

	fldErr = validateExecutor(spinv1.SpinAppSpec{Executor: constants.CyclotronExecutor}, func(string) bool { return false })
	require.EqualError(t, fldErr, "spec.executor: Invalid value: \"cyclotron\": executor does not exist on cluster")

	fldErr = validateExecutor(spinv1.SpinAppSpec{Executor: constants.ContainerDShimSpinExecutor}, func(name string) bool { return true })
	require.Nil(t, fldErr)
}

func TestValidateReplicas(t *testing.T) {
	t.Parallel()

	fldErr := validateReplicas(spinv1.SpinAppSpec{})
	require.EqualError(t, fldErr, "spec.replicas: Invalid value: 0: replicas must be > 0")

	fldErr = validateReplicas(spinv1.SpinAppSpec{Replicas: 1})
	require.Nil(t, fldErr)
}

func TestValidateAnnotations(t *testing.T) {
	t.Parallel()

	fldErr := validateAnnotations(spinv1.SpinAppSpec{
		Executor:              constants.CyclotronExecutor,
		DeploymentAnnotations: map[string]string{"key": "asdf"},
	})
	require.EqualError(t, fldErr,
		`spec.deploymentAnnotations: Invalid value: map[string]string{"key":"asdf"}: `+
			`deploymentAnnotations can't be set when runtime is cyclotron`)

	fldErr = validateAnnotations(spinv1.SpinAppSpec{
		Executor:       constants.CyclotronExecutor,
		PodAnnotations: map[string]string{"key": "asdf"},
	})
	require.EqualError(t, fldErr,
		`spec.podAnnotations: Invalid value: map[string]string{"key":"asdf"}: `+
			`podAnnotations can't be set when runtime is cyclotron`)

	fldErr = validateAnnotations(spinv1.SpinAppSpec{
		Executor:              constants.ContainerDShimSpinExecutor,
		DeploymentAnnotations: map[string]string{"key": "asdf"},
	})
	require.Nil(t, fldErr)

	fldErr = validateAnnotations(spinv1.SpinAppSpec{
		Executor: constants.CyclotronExecutor,
	})
	require.Nil(t, fldErr)
}
