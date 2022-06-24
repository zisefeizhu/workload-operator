package template

import (
	"github.com/zisefeizhu/workload-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewService 暂时命名
func NewService(workload *v1alpha1.Workload) *corev1.Service {
	if workload.Spec.ServiceType == "" {
		workload.Spec.ServiceType = corev1.ServiceTypeClusterIP
	}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workload.Name,
			Namespace: workload.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector:        workload.Labels,
			Type:            workload.Spec.ServiceType,
			Ports:           workload.Spec.ServicePorts,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
	if workload.Spec.HeadlessService {
		svc.Spec.ClusterIP = corev1.ClusterIPNone
		svc.Spec.Type = ""
	}
	return svc
}
