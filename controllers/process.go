package controllers

import (
	"context"
	"errors"
	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	"github.com/zisefeizhu/workload-operator/tool"
	appv1 "k8s.io/api/apps/v1"
)

// 处理wk status的 func
func (r *WorkloadReconciler) workloadStatus(instance *workloadsv1alpha1.Workload, dgStatus *workloadsv1alpha1.DeploymentGroupStatus, svcStatus *workloadsv1alpha1.ServiceStatus, ctx context.Context) (error, workloadsv1alpha1.WorkloadStatus) {
	// status
	s := workloadsv1alpha1.WorkloadStatus{}
	s.DeploymentGroupStatus = *dgStatus
	s.ServiceStatus = *svcStatus
	// 更改cr的状态
	// 计算工作负载和svc
	if *dgStatus.Replicas == 0 {
		s.Phase = workloadsv1alpha1.PendingPhase
	} else if *dgStatus.Replicas == dgStatus.AvailableReplicas {
		s.Phase = workloadsv1alpha1.RunningPhase
	} else if *dgStatus.Replicas != dgStatus.AvailableReplicas {
		s.Phase = workloadsv1alpha1.UpdatePhase
	}
	s.LastUpdateTime = tool.TimeToKubernetes()
	instance.Status = s
	// todo
	return r.Status().Update(ctx, instance), s
}

// 处理wk phase的 func
func (r *WorkloadReconciler) workloadPhase(ctx context.Context, instance *workloadsv1alpha1.Workload, phase workloadsv1alpha1.Phase) error {
	instance.Status.Phase = phase
	instance.Status.LastUpdateTime = tool.TimeToKubernetes()
	return r.Status().Update(ctx, instance)
}

func (r *WorkloadReconciler) workloadStatusProcessor(w interface{}) (resp *workloadsv1alpha1.DeploymentGroupStatus) {
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

// 工作负载矫正处理器
// 进入func 状态为非 pending和 running
func (r *WorkloadReconciler) workloadCorrectionProcessor(workloadStatus *workloadsv1alpha1.WorkloadStatus) error {
	switch workloadStatus.Phase {
	case workloadsv1alpha1.UnknownPhase:
		// todo
		break
	case workloadsv1alpha1.FailedPhase:
		// todo
		break
	case workloadsv1alpha1.UpdatePhase:
		// todo
		break
	default:
		return errors.New("工作负载矫正处理 传入参数错误")
	}
	return nil
}
