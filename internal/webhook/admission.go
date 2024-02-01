package webhook

import (
	spinv1 "github.com/spinkube/spin-operator/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func SetupSpinAppWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&spinv1.SpinApp{}).
		WithDefaulter(&SpinAppDefaulter{Client: mgr.GetClient()}).
		WithValidator(&SpinAppValidator{Client: mgr.GetClient()}).
		Complete()
}

func SetupSpinAppExecutorWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&spinv1.SpinAppExecutor{}).
		WithDefaulter(&SpinAppExecutorDefaulter{Client: mgr.GetClient()}).
		WithValidator(&SpinAppExecutorValidator{Client: mgr.GetClient()}).
		Complete()
}
