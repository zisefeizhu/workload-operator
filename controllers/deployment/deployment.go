package deployment

import (
	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deploymentClient struct {
	w *workloadsv1alpha1.Workload
}

func NewDeployment(w *workloadsv1alpha1.Workload) *deploymentClient {
	return &deploymentClient{
		w: w,
	}
}

func (d *deploymentClient) Create() interface{} {
	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.w.Name,
			Namespace: d.w.Namespace,
			Labels:    d.w.Labels,
		},
		Spec: appv1.DeploymentSpec{
			Replicas: d.w.Spec.Replicas,
			Selector: d.w.Spec.Selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: d.w.Labels,
				},
				Spec: d.w.Spec.Template.Spec,
			},
		},
	}
}

func (d *deploymentClient) Found() interface{} {
	return &appv1.Deployment{}
}
