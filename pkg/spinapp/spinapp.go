package spinapp

import (
	"fmt"

	"github.com/spinkube/spin-operator/internal/constants"
)

const (
	// HTTPPortName is the name used to identify the HTTP Port on a spin app
	// deployment.
	HTTPPortName = "http-app"

	// DefaultHTTPPort is the port that the operator will assign to a pod by
	// default when constructing deployments and services.
	DefaultHTTPPort = 80

	// StatusReady is the ready value for an app status label.
	StatusReady = "ready"
)

var (
	// NameLabelKey is the app name label key.
	NameLabelKey = constants.ConstructResourceLabelKey("app-name")
)

// ConstructStatusLabelKey returns the app status label key, used primarily
// in Service selectors.
func ConstructStatusLabelKey(appName string) string {
	return constants.ConstructResourceLabelKey(fmt.Sprintf("app.%s.status", appName))
}

// ConstructStatusReadyLabel returns the app status label key and value used
// by a Service selector to select Pod(s) ready to serve an app.
func ConstructStatusReadyLabel(appName string) (string, string) {
	return ConstructStatusLabelKey(appName), StatusReady
}
