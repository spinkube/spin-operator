package webhook

import (
	"context"

	spinv1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/internal/logging"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:lll
//+kubebuilder:webhook:path=/validate-core-spinoperator-dev-v1alpha1-spinappexecutor,mutating=false,failurePolicy=fail,sideEffects=None,groups=core.spinkube.dev,resources=spinappexecutors,verbs=create;update,versions=v1alpha1,name=vspinappexecutor.kb.io,admissionReviewVersions=v1

// SpinAppExecutorValidator validates SpinApps
type SpinAppExecutorValidator struct {
	Client client.Client
}

// ValidateCreate implements webhook.Validator
func (v *SpinAppExecutorValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	log := logging.FromContext(ctx)

	executor := obj.(*spinv1alpha1.SpinAppExecutor)
	log.Info("validate create", "name", executor.Name)

	return nil, v.validateSpinAppExecutor(executor)
}

// ValidateUpdate implements webhook.Validator
func (v *SpinAppExecutorValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	log := logging.FromContext(ctx)

	executor := newObj.(*spinv1alpha1.SpinAppExecutor)
	log.Info("validate update", "name", executor.Name)

	return nil, v.validateSpinAppExecutor(executor)
}

// ValidateDelete implements webhook.Validator
func (v *SpinAppExecutorValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	log := logging.FromContext(ctx)

	executor := obj.(*spinv1alpha1.SpinAppExecutor)
	log.Info("validate delete", "name", executor.Name)

	return nil, nil
}

func (v *SpinAppExecutorValidator) validateSpinAppExecutor(executor *spinv1alpha1.SpinAppExecutor) error {
	return nil
}
