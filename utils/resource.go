package utils

import (
	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	appv1 "k8s.io/api/apps/v1"
)

func WorkloadStatusProcessor(w interface{}) (resp *workloadsv1alpha1.DeploymentGroupStatus) {
	switch i := w.(type) {
	case *appv1.Deployment:
		resp = &workloadsv1alpha1.DeploymentGroupStatus{
			Type:                workloadsv1alpha1.Kind(i.Kind),
			AvailableReplicas:   i.Status.AvailableReplicas,
			Replicas:            i.Spec.Replicas,
			UnavailableReplicas: i.Status.UnavailableReplicas,
		}
	}
	return
}
