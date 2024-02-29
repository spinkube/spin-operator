package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type spinAppExecutorMetrics struct {
	// Record executor metadata as info metric.
	// This is particularly useful for querying
	// available executor names for creating dashboards
	infoGauge *prometheus.GaugeVec
}

func newSpinAppExecutorMetrics() *spinAppExecutorMetrics {
	return &spinAppExecutorMetrics{
		infoGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "spin_operator_spinapp_executor_info",
				Help: "info about spinapp executor labeled by name, namespace, create_deployment and runtime_class_name",
			},
			[]string{
				"name",
				"namespace",
				"create_deployment",
				"runtime_class_name",
			},
		),
	}
}

func (m *spinAppExecutorMetrics) Register(registry metrics.RegistererGatherer) {
	registry.MustRegister(m.infoGauge)
}

func (r *SpinAppExecutorReconciler) setupMetrics() {
	r.metrics = newSpinAppExecutorMetrics()
	r.metrics.Register(r.MetricsRegistry)
}
