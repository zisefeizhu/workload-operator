package controllers

import (
	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	"github.com/zisefeizhu/workload-operator/controllers/deployment"
)

type Workload interface {
	Create() interface{}
	Found() interface{}
}

func NewWorkload(w *workloadsv1alpha1.Workload) Workload {
	var cli Workload
	switch w.Spec.Type {
	case "deployment":
		cli = deployment.NewDeployment(w)
		break
	}
	return cli
}
