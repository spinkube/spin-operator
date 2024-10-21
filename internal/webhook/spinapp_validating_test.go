package webhook

import (
	"testing"

	spinv1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/stretchr/testify/require"
)

func TestValidateReplicas(t *testing.T) {
	t.Parallel()

	fldErr := validateReplicas(spinv1alpha1.SpinAppSpec{})
	require.EqualError(t, fldErr, "spec.replicas: Invalid value: 0: replicas must be > 0")

	fldErr = validateReplicas(spinv1alpha1.SpinAppSpec{Replicas: 1})
	require.Nil(t, fldErr)
}
