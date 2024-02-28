/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/prometheus/client_golang/prometheus"
	spinv1 "github.com/spinkube/spin-operator/api/v1"
	"github.com/spinkube/spin-operator/internal/logging"
)

// SpinAppExecutorReconciler reconciles a SpinAppExecutor object
type SpinAppExecutorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core.spinoperator.dev,resources=spinappexecutors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.spinoperator.dev,resources=spinappexecutors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.spinoperator.dev,resources=spinappexecutors/finalizers,verbs=update

// SetupWithManager sets up the controller with the Manager.
func (r *SpinAppExecutorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&spinv1.SpinAppExecutor{}).
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *SpinAppExecutorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx)
	log.Debug("Reconciling SpinAppExecutor")

	// Check if the SpinAppExecutor exists
	var executor spinv1.SpinAppExecutor
	if err := r.Client.Get(ctx, req.NamespacedName, &executor); err != nil {
		log.Error(err, "Unable to fetch SpinAppExecutor")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// SpinAppExecutor has been requested for deletion
	if !executor.DeletionTimestamp.IsZero() {
		err := r.handleDeletion(ctx, &executor)
		if err != nil {
			return ctrl.Result{}, err
		}

		err = r.removeFinalizer(ctx, &executor)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	runtimeClassName := ""
	if executor.Spec.DeploymentConfig != nil {
		runtimeClassName = executor.Spec.DeploymentConfig.RuntimeClassName
	}

	// record spin_operator_spinapp_executor_info metric
	spinOperatorSpinAppExecutorInfo.With(prometheus.Labels{
		"name":               executor.Name,
		"namespace":          executor.Namespace,
		"create_deployment":  fmt.Sprintf("%t", executor.Spec.CreateDeployment),
		"runtime_class_name": runtimeClassName,
	}).Set(1)

	// Make sure the finalizer is present
	err := r.ensureFinalizer(ctx, &executor)
	return ctrl.Result{}, client.IgnoreNotFound(err)
}

// handleDeletion makes sure no SpinApps are dependent on the SpinAppExecutor
// before allowing it to be deleted.
func (r *SpinAppExecutorReconciler) handleDeletion(ctx context.Context, executor *spinv1.SpinAppExecutor) error {
	var spinApps spinv1.SpinAppList
	if err := r.Client.List(ctx, &spinApps, client.MatchingFields{"spec.runtime": executor.Name}); err != nil {
		// TODO: Log this
		// TODO: Emit k8s event
		return err
	}

	if len(spinApps.Items) > 0 {
		return errors.New("cannot delete SpinAppExecutor with dependent SpinApps")
	}

	return nil
}

// removeFinalizer removes the finalizer from a SpinAppExecutor.
func (r *SpinAppExecutorReconciler) removeFinalizer(ctx context.Context, executor *spinv1.SpinAppExecutor) error {
	if controllerutil.ContainsFinalizer(executor, SpinOperatorFinalizer) {
		controllerutil.RemoveFinalizer(executor, SpinOperatorFinalizer)
		if err := r.Client.Update(ctx, executor); err != nil {
			return err
		}
	}
	return nil
}

// ensureFinalizer ensures the finalizer is present on a SpinAppExecutor.
func (r *SpinAppExecutorReconciler) ensureFinalizer(ctx context.Context, executor *spinv1.SpinAppExecutor) error {
	if !controllerutil.ContainsFinalizer(executor, SpinOperatorFinalizer) {
		controllerutil.AddFinalizer(executor, SpinOperatorFinalizer)
		if err := r.Client.Update(ctx, executor); err != nil {
			return err
		}
	}
	return nil
}
