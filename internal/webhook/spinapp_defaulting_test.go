package webhook

import (
	"context"
	"testing"

	spinv1 "github.com/spinkube/spin-operator/api/v1"
	"github.com/spinkube/spin-operator/internal/constants"
	"github.com/stretchr/testify/require"
)

func TestDefaultNothingToSet(t *testing.T) {
	t.Parallel()

	defaulter := &SpinAppDefaulter{}

	spinApp := &spinv1.SpinApp{Spec: spinv1.SpinAppSpec{
		Executor: constants.CyclotronExecutor,
		Replicas: 1,
	}}

	err := defaulter.Default(context.Background(), spinApp)
	require.NoError(t, err)
	require.Equal(t, constants.CyclotronExecutor, spinApp.Spec.Executor)
	require.Equal(t, int32(1), spinApp.Spec.Replicas)
}
