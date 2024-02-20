package webhook

import (
	"fmt"
	"net/http"
	"os"
	"time"

	spinv1 "github.com/spinkube/spin-operator/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func LazyWebhookStarter(mgr ctrl.Manager) error {
	ticker := time.NewTicker(2 * time.Second)
	timeout := time.NewTimer(5 * time.Minute)

	crtFile := "/tmp/k8s-webhook-server/serving-certs/tls.crt"

	for {
		select {
		case <-ticker.C:
			_, err := os.ReadFile(crtFile)
			if err != nil && os.IsNotExist(err) {
				fmt.Printf("file %s does not exist yet\n", crtFile)
				continue
			}

			fmt.Printf("crtfile found, setting up webhook")

			webhookSetupLog := ctrl.Log.WithName("webhook-setup")
			if err = SetupSpinAppWebhookWithManager(mgr); err != nil {
				webhookSetupLog.Error(err, "unable to create webhook", "webhook", "SpinApp")
				os.Exit(1)
			}
			if err = SetupSpinAppExecutorWebhookWithManager(mgr); err != nil {
				webhookSetupLog.Error(err, "unable to create webhook", "webhook", "SpinAppExecutor")
				os.Exit(1)
			}

			mgr.GetWebhookServer().WebhookMux().HandleFunc("webhooks-ready", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})

			return nil
		case <-timeout.C:
			ticker.Stop()
			panic("timed out while waiting for webhook to start")
		}
	}
}

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
