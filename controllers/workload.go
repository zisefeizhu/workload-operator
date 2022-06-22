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
	case "deployment":
		cli = deployment.NewDeployment(w)
		break
	case "statefulSet":
		cli = statefulSet.NewStatefulSet(w)
		break
	case "daemonSet":
		cli = daemonSet.NewDaemonSet(w)
		break
	case "job":
		cli = job.NewJob(w)
		break
	case "cronJob":
		cli = cronJob.NewCronJob(w)
		break
	}
	return cli
}
