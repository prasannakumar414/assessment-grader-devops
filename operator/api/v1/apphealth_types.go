/*
Copyright 2026.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppHealthSpec defines the desired state of AppHealth (what to validate).
type AppHealthSpec struct {
	// Namespace where the student deployed their app.
	Namespace string `json:"namespace"`
	// AppName is the name of the Deployment and Service the student should have created.
	AppName string `json:"appName"`
	// ExpectedImage (optional): if set, validates that the deployment uses this image.
	ExpectedImage string `json:"expectedImage,omitempty"`
	// InfoPath is the path to call for email (default: /info).
	InfoPath string `json:"infoPath,omitempty"`
	// ReportURL (optional): external API to report validation results.
	ReportURL string `json:"reportURL,omitempty"`
}

// AppHealthStatus defines the observed state of AppHealth.
type AppHealthStatus struct {
	// DeploymentStatus: Healthy if deployment exists and has ready replicas.
	DeploymentStatus string `json:"deploymentStatus,omitempty"`
	// ServiceStatus: Healthy if service exists and has ports.
	ServiceStatus string `json:"serviceStatus,omitempty"`
	// ImageValid is true if ExpectedImage is not set or matches the deployed image.
	ImageValid bool `json:"imageValid,omitempty"`
	// Healthy is true when deployment, service, and image are all valid.
	Healthy bool `json:"healthy,omitempty"`
	// Email from the student app's /info endpoint.
	Email string `json:"email,omitempty"`
	// LastChecked timestamp.
	LastChecked string `json:"lastChecked,omitempty"`
	// Message for grading feedback.
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AppHealth is the Schema for the apphealths API (student deployment validation).
type AppHealth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppHealthSpec   `json:"spec,omitempty"`
	Status AppHealthStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AppHealthList contains a list of AppHealth.
type AppHealthList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppHealth `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppHealth{}, &AppHealthList{})
}
