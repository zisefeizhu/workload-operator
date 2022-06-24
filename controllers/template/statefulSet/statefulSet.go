package statefulSet

import (
	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type statefulSetClient struct {
	w *workloadsv1alpha1.Workload
}

func NewStatefulSet(w *workloadsv1alpha1.Workload) *statefulSetClient {
	return &statefulSetClient{
		w: w,
	}
}

const kind = "StatefulSet"

func (s *statefulSetClient) Template() interface{} {
	return &appv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind: kind,
			//APIVersion: aPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.w.Name,
			Namespace: s.w.Namespace,
			Labels:    s.w.Labels,
		},
		Spec: appv1.StatefulSetSpec{
			ServiceName: s.w.Name,
			Selector:    s.w.Spec.WorkloadSpec.Selector,
			Replicas:    s.w.Spec.WorkloadSpec.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: s.w.Labels,
				},
				Spec: s.w.Spec.WorkloadSpec.Template.Spec,
			},
			VolumeClaimTemplates: s.w.Spec.WorkloadSpec.VolumeClaimTemplates,
		},
	}
}

func (s *statefulSetClient) Found() interface{} {
	return &appv1.StatefulSet{}
}
