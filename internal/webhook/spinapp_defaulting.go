package webhook

import (
	"context"
	"strings"

	spinv1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/internal/logging"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// nolint:lll
//+kubebuilder:webhook:path=/mutate-core-spinoperator-dev-v1alpha1-spinapp,mutating=true,failurePolicy=fail,sideEffects=None,groups=core.spinoperator.dev,resources=spinapps,verbs=create;update,versions=v1alpha1,name=mspinapp.kb.io,admissionReviewVersions=v1

// SpinAppDefaulter mutates SpinApps
type SpinAppDefaulter struct {
	Client client.Client
}

// Default implements webhook.Defaulter
func (d *SpinAppDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	log := logging.FromContext(ctx)

	spinApp := obj.(*spinv1alpha1.SpinApp)
	log.Info("default", "name", spinApp.Name)

	if spinApp.Spec.Executor == "" {
		executor, err := d.findDefaultExecutor(ctx)
		if err != nil {
			return err
		}
		spinApp.Spec.Executor = executor
	}

	return nil
}

// findDefaultExecutor sets the default executor for a SpinApp.
//
// Defaults to whatever executor is available on the cluster. If multiple
// executors are available then the first executor in alphabetical order
// will be chosen. If no executors are available then no default will be set.
func (d *SpinAppDefaulter) findDefaultExecutor(ctx context.Context) (string, error) {
	log := logging.FromContext(ctx)

	var executors spinv1alpha1.SpinAppExecutorList
	if err := d.Client.List(ctx, &executors); err != nil {
		log.Error(err, "failed to list SpinAppExecutors")
		return "", err
	}

	if len(executors.Items) == 0 {
		log.Info("no SpinAppExecutors found")
		return "", nil
	}

	// Return first executor in alphabetical order
	chosenExecutor := executors.Items[0]
	for _, executor := range executors.Items[1:] {
		// For each item after the first see if it is alphabetically before the current chosen executor
		if strings.Compare(executor.Name, chosenExecutor.Name) < 0 {
			chosenExecutor = executor
		}
	}

	log.Info("defaulting to executor", "name", chosenExecutor.Name)
	return chosenExecutor.Name, nil
}
