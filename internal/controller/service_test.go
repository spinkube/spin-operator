package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstructService(t *testing.T) {
	t.Parallel()

	app := minimalSpinApp()
	svc := constructService(app)

	// Omitting these is likely to result in client breakage, thus require a test
	// change.
	require.Equal(t, "Service", svc.TypeMeta.Kind)
	require.Equal(t, "v1", svc.TypeMeta.APIVersion)

	// We expect that the service object has the app name and nothing else.
	require.Equal(t, map[string]string{"core.spinoperator.dev/app-name": "my-app"}, svc.ObjectMeta.Labels)
	// We expect that the service selector has the app status and nothing else.
	require.Equal(t, map[string]string{"core.spinoperator.dev/app.my-app.status": "ready"}, svc.Spec.Selector)

	// We expect that the HTTP Port is part of the service. There's currently no
	// non-http implementations of a Spin trigger in Kubernetes, thus nothing that
	// would change this.
	require.Len(t, svc.Spec.Ports, 1)
	require.Equal(t, int32(80), svc.Spec.Ports[0].Port)
	require.Equal(t, "http-app", svc.Spec.Ports[0].TargetPort.StrVal)
}
