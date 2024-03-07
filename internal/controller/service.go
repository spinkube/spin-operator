package controller

import (
	spinv1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/pkg/spinapp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// constructService builds a corev1.Service based on the configuration of a SpinApp.
func constructService(app *spinv1alpha1.SpinApp) *corev1.Service {
	annotations := app.Spec.ServiceAnnotations
	if annotations == nil {
		annotations = map[string]string{}
	}

	labels := constructAppLabels(app)

	statusKey, statusValue := spinapp.ConstructStatusReadyLabel(app.Name)
	selector := map[string]string{statusKey: statusValue}

	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        app.Name,
			Namespace:   app.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						Type:   intstr.String,
						StrVal: spinapp.HTTPPortName,
					},
					Port: spinapp.DefaultHTTPPort,
				},
			},
			Selector: selector,
		},
	}

	return svc
}

// constructAppLabels returns the labels to add to deployment/service
// objects for the given SpinApp
func constructAppLabels(app *spinv1alpha1.SpinApp) map[string]string {
	return map[string]string{
		spinapp.NameLabelKey: app.Name,
	}
}
