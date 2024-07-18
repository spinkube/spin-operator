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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SpinAppExecutorSpec defines the desired state of SpinAppExecutor
type SpinAppExecutorSpec struct {
	// CreateDeployment specifies whether the Executor wants the SpinKube operator
	// to create a deployment for the application or if it will be realized externally.
	CreateDeployment bool `json:"createDeployment"`

	// DeploymentConfig specifies how the deployment should be configured when
	// createDeployment is true.
	DeploymentConfig *ExecutorDeploymentConfig `json:"deploymentConfig,omitempty"`
}

type ExecutorDeploymentConfig struct {
	// RuntimeClassName is the runtime class name that should be used by pods created
	// as part of a deployment. This should only be defined when SpintainerImage is not defined.
	RuntimeClassName *string `json:"runtimeClassName,omitempty"`

	// SpinImage points to an image that will run Spin in a container to execute
	// your SpinApp. This is an alternative to using the shim to execute your
	// SpinApp. This should only be defined when RuntimeClassName is not
	// defined. When specified, application images must be available without
	// authentication.
	SpinImage *string `json:"spinImage,omitempty"`

	// CACertSecret specifies the name of the secret containing the CA
	// certificates to be mounted to the deployment.
	CACertSecret string `json:"caCertSecret,omitempty"`

	// InstallDefaultCACerts specifies whether the default CA
	// certificate bundle should be generated. When set a new secret
	// will be created containing the certificates. If no secret name is
	// defined in `CACertSecret` the secret name will be `spin-ca`.
	InstallDefaultCACerts bool `json:"installDefaultCACerts,omitempty"`

	// Otel provides Kubernetes Bindings to Otel Variables.
	Otel *OtelConfig `json:"otel,omitempty"`
}

// SpinAppExecutorStatus defines the observed state of SpinAppExecutor
type SpinAppExecutorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// OtelConfig is the supported environment variables for OpenTelemetry
type OtelConfig struct {
	// ExporterOtlpEndpoint configures the default combined otlp endpoint for sending telemetry
	ExporterOtlpEndpoint string `json:"exporter_otlp_endpoint,omitempty"`
	// ExporterOtlpTracesEndpoint configures the trace-specific otlp endpoint
	ExporterOtlpTracesEndpoint string `json:"exporter_otlp_traces_endpoint,omitempty"`
	// ExporterOtlpMetricsEndpoint configures the metrics-specific otlp endpoint
	ExporterOtlpMetricsEndpoint string `json:"exporter_otlp_metrics_endpoint,omitempty"`
	// ExporterOtlpLogsEndpoint configures the logs-specific otlp endpoint
	ExporterOtlpLogsEndpoint string `json:"exporter_otlp_logs_endpoint,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SpinAppExecutor is the Schema for the spinappexecutors API
type SpinAppExecutor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpinAppExecutorSpec   `json:"spec,omitempty"`
	Status SpinAppExecutorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SpinAppExecutorList contains a list of SpinAppExecutor
type SpinAppExecutorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SpinAppExecutor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SpinAppExecutor{}, &SpinAppExecutorList{})
}
