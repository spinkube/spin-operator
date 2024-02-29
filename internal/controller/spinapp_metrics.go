package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type spinAppMetrics struct {
	infoGauge *prometheus.GaugeVec
}

func newSpinAppMetrics() *spinAppMetrics {
	return &spinAppMetrics{
		infoGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "spin_operator_spinapp_info",
				Help: "info about spinapp labeled by name, namespace, executor",
			},
			[]string{
				"name",
				"namespace",
				"executor",
			},
		),
	}
}

func (m *spinAppMetrics) Register(registry metrics.RegistererGatherer) {
	registry.MustRegister(m.infoGauge)
}

func (r *SpinAppReconciler) setupMetrics() {
	r.metrics = newSpinAppMetrics()
	r.metrics.Register(r.MetricsRegistry)
}
