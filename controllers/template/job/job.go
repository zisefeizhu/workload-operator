package job

import (
	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type jobClient struct {
	w *workloadsv1alpha1.Workload
}

func NewJob(w *workloadsv1alpha1.Workload) *jobClient {
	return &jobClient{
		w: w,
	}
}

const kind = "Job"

func (j *jobClient) Template() interface{} {
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind: kind,
			//APIVersion: aPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.w.Name,
			Namespace: j.w.Namespace,
			Labels:    j.w.Labels,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: j.w.Spec.WorkloadSpec.JobBackoffLimit,
			//Selector:     j.w.Spec.Selector,    job selector 不允许自定义,controller  生成 controller-uid
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: j.w.Labels,
				},
				Spec: j.w.Spec.WorkloadSpec.Template.Spec,
			},
		},
	}
}

func (j *jobClient) Found() interface{} {
	// job
	return &batchv1.Job{}
}
