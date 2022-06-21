/*
Copyright 2022.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WorkloadSpec defines the desired state of Workload
type WorkloadSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Workload. Edit workload_types.go to remove/update
	//Foo string `json:"foo,omitempty"`
	// 部署组类型
	Type string `json:"type"`
	// 副本数
	Replicas *int32 `json:"replicas,omitempty"`
	// 是否启用service
	EnableService bool `json:"enableService,omitempty"`
	// 标签选择器
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
	// pod模版
	Template *corev1.PodTemplateSpec `json:"template"`
	// service 类型
	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`
	// service 端口
	ServicePorts []corev1.ServicePort `json:"servicePorts,omitempty"`
}

//type workloadService struct {
//	// service 类型
//	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`
//	// service 端口
//	ServicePorts []corev1.ServicePort `json:"servicePorts,omitempty"`
//}

// WorkloadStatus defines the observed state of Workload
type WorkloadStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:printcolumn:JSONPath=".spec.type",name=Type,type=string
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Workload is the Schema for the workloads API
type Workload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WorkloadSpec   `json:"spec,omitempty"`
	Status            WorkloadStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WorkloadList contains a list of Workload
type WorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workload `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Workload{}, &WorkloadList{})
}
