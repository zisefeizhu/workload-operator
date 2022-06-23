package controllers

import (
	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	"github.com/zisefeizhu/workload-operator/controllers/template/cronJob"
	"github.com/zisefeizhu/workload-operator/controllers/template/daemonSet"
	"github.com/zisefeizhu/workload-operator/controllers/template/deployment"
	"github.com/zisefeizhu/workload-operator/controllers/template/job"
	"github.com/zisefeizhu/workload-operator/controllers/template/statefulSet"
)

type Workload interface {
	Template() interface{}
	Found() interface{}
}

func NewWorkload(w *workloadsv1alpha1.Workload) Workload {
	var cli Workload
	switch w.Spec.Type {
	case workloadsv1alpha1.DeploymentKind:
		cli = deployment.NewDeployment(w)
		break
	case workloadsv1alpha1.StatefulSetKind:
		cli = statefulSet.NewStatefulSet(w)
		break
	case workloadsv1alpha1.DaemonSetKind:
		cli = daemonSet.NewDaemonSet(w)
		break
	case workloadsv1alpha1.JobKind:
		cli = job.NewJob(w)
		break
	case workloadsv1alpha1.CronjobKind:
		cli = cronJob.NewCronJob(w)
		break
	}
	return cli
}
