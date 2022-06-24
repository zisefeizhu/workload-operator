package daemonSet

import (
	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type daemonSetClient struct {
	w *workloadsv1alpha1.Workload
}

func NewDaemonSet(w *workloadsv1alpha1.Workload) *daemonSetClient {
	return &daemonSetClient{
		w: w,
	}
}

func (d *daemonSetClient) Template() interface{} {
	return &appv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.w.Name,
			Namespace: d.w.Namespace,
			Labels:    d.w.Labels,
		},
		Spec: appv1.DaemonSetSpec{
			Selector: d.w.Spec.WorkloadSpec.Selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: d.w.Labels,
				},
				Spec: d.w.Spec.WorkloadSpec.Template.Spec,
			},
		},
	}
}

func (d *daemonSetClient) Found() interface{} {
	return &appv1.DaemonSet{}
}
