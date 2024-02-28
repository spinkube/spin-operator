package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

func init() {
	registerPrometheusMetrics()
}

func registerPrometheusMetrics() {
	metrics.Registry.MustRegister(spinOperatorSpinAppInfo)
	metrics.Registry.MustRegister(spinOperatorSpinAppExecutorInfo)
}

var (
	spinOperatorSpinAppInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "spin_operator_spinapp_info",
			Help: "info about spinapp labeled by name, namespace, executor",
		},
		[]string{
			"name",
			"namespace",
			"executor",
		},
	)

	spinOperatorSpinAppExecutorInfo = prometheus.NewGaugeVec(
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
	)
)
