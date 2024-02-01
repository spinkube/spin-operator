package webhook

import (
	"context"

	spinv1 "github.com/spinkube/spin-operator/api/v1"
	"github.com/spinkube/spin-operator/internal/constants"
	"github.com/spinkube/spin-operator/internal/logging"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:lll
//+kubebuilder:webhook:path=/validate-core-spinoperator-dev-v1-spinapp,mutating=false,failurePolicy=fail,sideEffects=None,groups=core.spinoperator.dev,resources=spinapps,verbs=create;update,versions=v1,name=vspinapp.kb.io,admissionReviewVersions=v1

// SpinAppValidator validates SpinApps
type SpinAppValidator struct {
	Client client.Client
}

// ValidateCreate implements webhook.Validator
func (v *SpinAppValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	log := logging.FromContext(ctx)

	spinApp := obj.(*spinv1.SpinApp)
	log.Info("validate create", "name", spinApp.Name)

	return nil, v.validateSpinApp(ctx, spinApp)
}

// ValidateUpdate implements webhook.Validator
func (v *SpinAppValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	log := logging.FromContext(ctx)

	spinApp := newObj.(*spinv1.SpinApp)
	log.Info("validate update", "name", spinApp.Name)

	return nil, v.validateSpinApp(ctx, spinApp)
}

// ValidateDelete implements webhook.Validator
func (v *SpinAppValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	log := logging.FromContext(ctx)

	spinApp := obj.(*spinv1.SpinApp)
	log.Info("validate delete", "name", spinApp.Name)

	return nil, nil
}

func (v *SpinAppValidator) validateSpinApp(ctx context.Context, spinApp *spinv1.SpinApp) error {
	var allErrs field.ErrorList
	if err := validateExecutor(spinApp.Spec, v.executorExists(ctx, spinApp.Namespace)); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := validateReplicas(spinApp.Spec); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := validateAnnotations(spinApp.Spec); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "core.spinoperator.dev", Kind: "SpinApp"},
		spinApp.Name, allErrs)
}

// executorExists returns a function that checks if an executor exists on the cluster.
//
// We assume that the executor must exist in the same namespace as the SpinApp.
func (v *SpinAppValidator) executorExists(ctx context.Context, spinAppNs string) func(string) bool {
	return func(name string) bool {
		var executor spinv1.SpinAppExecutor
		if err := v.Client.Get(ctx, client.ObjectKey{Name: name, Namespace: spinAppNs}, &executor); err != nil {
			// TODO: This groups in both not found and client errors. We should ideally separate the two.
			return false
		}

		return true
	}
}

func validateExecutor(spec spinv1.SpinAppSpec, executorExists func(name string) bool) *field.Error {
	if spec.Executor == "" {
		return field.Invalid(field.NewPath("spec").Child("executor"), spec.Executor, "executor must be set, likely no default executor was set because you have no executors installed")
	}
	if !executorExists(spec.Executor) {
		return field.Invalid(field.NewPath("spec").Child("executor"), spec.Executor, "executor does not exist on cluster")
	}

	return nil
}

func validateReplicas(spec spinv1.SpinAppSpec) *field.Error {
	if spec.EnableAutoscaling && spec.Replicas != 0 {
		return field.Invalid(field.NewPath("spec").Child("replicas"), spec.Replicas, "replicas cannot be set when autoscaling is enabled")
	}
	if !spec.EnableAutoscaling && spec.Replicas < 1 {
		return field.Invalid(field.NewPath("spec").Child("replicas"), spec.Replicas, "replicas must be > 0")
	}

	return nil
}

func validateAnnotations(spec spinv1.SpinAppSpec) *field.Error {
	if spec.Executor != constants.CyclotronExecutor {
		return nil
	}
	if len(spec.DeploymentAnnotations) != 0 {
		return field.Invalid(field.NewPath("spec").Child("deploymentAnnotations"), spec.DeploymentAnnotations, "deploymentAnnotations can't be set when runtime is cyclotron")
	}
	if len(spec.PodAnnotations) != 0 {
		return field.Invalid(field.NewPath("spec").Child("podAnnotations"), spec.PodAnnotations, "podAnnotations can't be set when runtime is cyclotron")
	}

	return nil
}
