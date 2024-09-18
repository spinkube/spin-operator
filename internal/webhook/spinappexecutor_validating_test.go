package webhook

import (
	"testing"

	spinv1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/internal/generics"
	"github.com/stretchr/testify/require"
)

func TestValidateRuntimeClassAndSpinImage(t *testing.T) {
	t.Parallel()

	fldErr := validateRuntimeClassAndSpinImage(&spinv1alpha1.SpinAppExecutorSpec{
		CreateDeployment: true,
		DeploymentConfig: &spinv1alpha1.ExecutorDeploymentConfig{
			RuntimeClassName: generics.Ptr("foo"),
			SpinImage:        generics.Ptr("bar"),
		},
	})
	require.EqualError(t, fldErr, "spec.deploymentConfig.runtimeClassName: Invalid value: \"foo\": runtimeClassName and spinImage are mutually exclusive")

	fldErr = validateRuntimeClassAndSpinImage(&spinv1alpha1.SpinAppExecutorSpec{
		CreateDeployment: true,
		DeploymentConfig: &spinv1alpha1.ExecutorDeploymentConfig{
			RuntimeClassName: generics.Ptr("foo"),
			SpinImage:        nil,
		},
	})
	require.Nil(t, fldErr)

	fldErr = validateRuntimeClassAndSpinImage(&spinv1alpha1.SpinAppExecutorSpec{
		CreateDeployment: true,
		DeploymentConfig: &spinv1alpha1.ExecutorDeploymentConfig{
			RuntimeClassName: nil,
			SpinImage:        generics.Ptr("bar"),
		},
	})
	require.Nil(t, fldErr)

	fldErr = validateRuntimeClassAndSpinImage(&spinv1alpha1.SpinAppExecutorSpec{
		CreateDeployment: true,
		DeploymentConfig: &spinv1alpha1.ExecutorDeploymentConfig{
			RuntimeClassName: nil,
			SpinImage:        nil,
		},
	})
	require.EqualError(t, fldErr, "spec.deploymentConfig.runtimeClassName: Invalid value: \"null\": either runtimeClassName or spinImage must be set")
}
