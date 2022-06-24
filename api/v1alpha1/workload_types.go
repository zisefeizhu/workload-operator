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
	"strings"
)

type Kind string

const (
	DeploymentKind  Kind = "deployment"
	StatefulSetKind Kind = "statefulSet"
	DaemonSetKind   Kind = "daemonSet"
	CronjobKind     Kind = "cronjob"
	JobKind         Kind = "job"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// WorkloadSpec defines the desired state of Workload
type WorkloadSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// 部署组类型
	//+kubebuilder:validation:Enum="deployment";"statefulSet";"daemonSet";"job";"cronJob"
	Type Kind `json:"type"`
	// 副本数
	Replicas *int32 `json:"replicas,omitempty"`
	// 标签选择器
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
	// pod模版
	Template *corev1.PodTemplateSpec `json:"template"`
	// statefulSet 存储模版
	VolumeClaimTemplates []corev1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
	// job 重试次数
	//backoffLimit
	JobBackoffLimit *int32 `json:"jobBackoffLimit,omitempty"`
	// crontabJob
	// schedule
	Schedule string `json:"schedule,omitempty"`
}

type SvcSpec struct {
	// 是否启用service
	EnableService bool `json:"enableService,omitempty"`
	// service 类型
	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`
	// service 端口
	ServicePorts []corev1.ServicePort `json:"servicePorts,omitempty"`
	// statefulSet 无头服务
	//clusterIP: None
	HeadlessService bool `json:"headlessService,omitempty"`
}

type Phase string

const (
	RunningPhase Phase = "Running"
	UpdatePhase  Phase = "Update"
	PendingPhase Phase = "Pending"
	UnknownPhase Phase = "Unknown"
	FailedPhase  Phase = "Failed"
)

// WorkloadStatus defines the observed state of Workload
type WorkloadStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// workloads 状态
	// 要不要设置默认状态？？？
	////+kubebuilder:default:="Unknown"
	Phase                 Phase `json:"phase,omitempty"`
	DeploymentGroupStatus `json:"deploymentGroupStatus,omitempty"`
	ServiceStatus         `json:"serviceStatus,omitempty"`
}

// 需要将工作负载的状态等信息给返出来的，和workload的status字段进行对比用
type DeploymentGroupStatus struct {
	// 返回工作负载的类型
	Type Kind `json:"type" binding:"required"`
	// AvailableReplicas 可用副本数
	AvailableReplicas int32 `json:"availableReplicas"`
	// replicas  期望副本数
	Replicas *int32 `json:"replicas,omitempty"`
	// UnavailableReplicas 不可用副本数
	UnavailableReplicas int32 `json:"unavailableReplicas,omitempty"`
}

type ServiceStatus struct {
	// serviceIP  如果EnableService = true 则输出serviceIP
	ServiceIP string `json:"serviceIP,omitempty"`
}

//+kubebuilder:printcolumn:JSONPath=".status.deploymentGroupStatus.type",name=Type,type=string
//+kubebuilder:printcolumn:JSONPath=".status.phase",name=Phase,type=string
//+kubebuilder:printcolumn:JSONPath=".status.deploymentGroupStatus.replicas",name=Replicas,type=string
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=wk

// Workload is the Schema for the workloads API
type Workload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Spec           `json:"spec,omitempty"`
	Status            WorkloadStatus `json:"status,omitempty"`
}

type Spec struct {
	WorkloadSpec WorkloadSpec `json:"workloadSpec"`
	SvcSpec      SvcSpec      `json:"svcSpec,omitempty"`
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

// wk的状态
// 1、如果EnableService = true ，则需要总和判断svc和workload的状态1
// 2、如果EnableService = false，则需要根据期望副本数和可用副本数综合判断
// type 种类转大写
func (w *Workload) StatusTypeToUpper(k Kind) Kind {
	s := string(k)
	return Kind(strings.ToUpper(s[:1]) + s[1:])
}

//func (w *WorkloadStatus) workloadPhase() string {
//	return ""
//}
