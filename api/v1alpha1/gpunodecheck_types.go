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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GpuNodeCheckSpec defines the desired state of GpuNodeCheck
type GpuNodeCheckSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// NodeSelector selects candidate GPU nodes.
	// Example: {"accelerator": "nvidia"}
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// RequiredLabels are labels that every selected GPU node should have.
	// Each entry is a label key (values are not checked).
	// +optional
	RequiredLabels []string `json:"requiredLabels,omitempty"`
}

// GpuNodeCheckStatus defines the observed state of GpuNodeCheck.
type GpuNodeCheckStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// TotalNodes is the number of selected candidate nodes.
	// +optional
	TotalNodes int32 `json:"totalNodes,omitempty"`

	// ReadyNodes is the number of selected nodes that are Ready.
	// +optional
	ReadyNodes int32 `json:"readyNodes,omitempty"`

	// MissingLabelNodes contains nodes missing required labels.
	// Each entry is a node name. (Doesn't indicate which label was missing.)
	// +optional
	MissingLabelNodes []string `json:"missingLabelNodes,omitempty"`

	// LastCheckedTime is the last reconciliation timestamp.
	// +optional
	LastCheckedTime *metav1.Time `json:"lastCheckedTime,omitempty"`

	// Conditions describe the qualification result.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Total",type=integer,JSONPath=`.status.totalNodes`
// +kubebuilder:printcolumn:name="Ready",type=integer,JSONPath=`.status.readyNodes`
// +kubebuilder:printcolumn:name="LastChecked",type=date,JSONPath=`.status.lastCheckedTime`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// GpuNodeCheck is the Schema for the gpunodechecks API
type GpuNodeCheck struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of GpuNodeCheck
	// +required
	Spec GpuNodeCheckSpec `json:"spec"`

	// status defines the observed state of GpuNodeCheck
	// +optional
	Status GpuNodeCheckStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// GpuNodeCheckList contains a list of GpuNodeCheck
type GpuNodeCheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []GpuNodeCheck `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GpuNodeCheck{}, &GpuNodeCheckList{})
}
