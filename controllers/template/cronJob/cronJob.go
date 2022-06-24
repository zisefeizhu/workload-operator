package cronJob

import (
	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type cronJobClient struct {
	w *workloadsv1alpha1.Workload
}

func NewCronJob(w *workloadsv1alpha1.Workload) *cronJobClient {
	return &cronJobClient{
		w: w,
	}
}

//todo
/*
   现kubebuilder 的版本为3.1 默认的k8s版本为1.19 ，crontabJob的api 为batchv1beta1
   在k8s 1.21 crontabJob的版本稳定为为batchv1
   考虑升级kubebuilder版本3.1 --> 3.4
*/

func (c *cronJobClient) Template() interface{} {
	return &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.w.Name,
			Namespace: c.w.Namespace,
			Labels:    c.w.Labels,
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule: c.w.Spec.WorkloadSpec.Schedule,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: c.w.Spec.WorkloadSpec.Template.Spec,
					},
				},
			},
		},
	}
}

func (c *cronJobClient) Found() interface{} {
	return &batchv1beta1.CronJob{}
}
