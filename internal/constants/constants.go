package constants

import "fmt"

// OperatorResourceKeyspace is the keyspace used for constructing application
// metadata on Kubernetes objects
const OperatorResourceKeyspace = "core.spinoperator.dev"

// ConstructResourceLabelKey is used when building operator-managed labels for
// resources.
func ConstructResourceLabelKey(kind string) string {
	return fmt.Sprintf("%s/%s", OperatorResourceKeyspace, kind)
}

// KnownExecutor is an enumeration of the executors that are well-known and
// supported by the spin operator.
type KnownExecutor string

const (
	ContainerDShimSpinExecutor = "containerd-shim-spin"
	CyclotronExecutor          = "cyclotron"
	SpinInContainer            = "spin-in-container"
)
