package controllers

import (
	"context"
	"fmt"
	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	"github.com/zisefeizhu/workload-operator/controllers/template"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// 集中处理 部署组
func (r *WorkloadReconciler) deploymentGroup(instance *workloadsv1alpha1.Workload, ctx context.Context, req ctrl.Request) (*workloadsv1alpha1.DeploymentGroupStatus, error) {
	// 获取工作负载模版
	w := template.NewWorkload(instance).Template()
	if err := controllerutil.SetControllerReference(instance, w.(metav1.Object), r.Scheme); err != nil {
		return nil, err
	}
	found := template.NewWorkload(instance).Found()
	if err := r.Get(ctx, req.NamespacedName, found.(client.Object)); err != nil {
		if errors.IsNotFound(err) {
			err := r.Create(ctx, w.(client.Object))
			if err != nil {
				r.Logger.Error(err, "create app failed")
				return nil, err
			}
			r.Recorder.Event(instance, corev1.EventTypeNormal, fmt.Sprintf("%s-controller", instance.Spec.WorkloadSpec.Type), fmt.Sprintf("type is %s name is %s create in  %s namespace", instance.Spec.WorkloadSpec.Type, instance.Name, instance.Namespace))
		}
		if !errors.IsNotFound(err) {
			return nil, err
		}
	} else {
		if err := r.Update(ctx, w.(client.Object)); err != nil {
			r.Logger.Error(err, "update app failed")
			return nil, err
		}
		//r.Recorder.Event(instance, corev1.EventTypeNormal, fmt.Sprintf("%s-controller", instance.Spec.Type), fmt.Sprintf("type is %s name is %s  update in %s namespace", instance.Spec.Type, instance.Name, instance.Namespace))
	}
	// todo
	// 处理工作负载的status
	// 类型断言
	return workloadStatusProcessor(w), nil
}

// 处理svc的 func
func (r *WorkloadReconciler) svc(instance *workloadsv1alpha1.Workload, ctx context.Context) (*workloadsv1alpha1.ServiceStatus, error) {
	// service的处理
	service := template.NewService(instance)
	err := controllerutil.SetControllerReference(instance, service, r.Scheme)
	if err != nil {
		return nil, err
	}
	s := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, s); err != nil {
		if errors.IsNotFound(err) && instance.Spec.SvcSpec.EnableService {
			if err := r.Create(ctx, service); err != nil {
				r.Logger.Error(err, "create service failed")
				return nil, err
			}
			r.Recorder.Event(instance, corev1.EventTypeNormal, fmt.Sprintf("%s-controller", "service"), fmt.Sprintf("type is %s name is %s create in  %s namespace", "serivce", instance.Name, instance.Namespace))
		}
		if !errors.IsNotFound(err) && instance.Spec.SvcSpec.EnableService {
			return nil, err
		}
	} else {
		if instance.Spec.SvcSpec.EnableService {
			currIP := s.Spec.ClusterIP
			s.Spec.ClusterIP = currIP
			service.Spec.ClusterIP = currIP
			if !reflect.DeepEqual(s.Spec, service.Spec) {
				s.Spec = service.Spec
				if err := r.Update(ctx, s); err != nil {
					r.Logger.Error(err, "update service failed")
					return nil, err
				}
				r.Recorder.Event(instance, corev1.EventTypeNormal, fmt.Sprintf("%s-controller", "service"), fmt.Sprintf("type is %s name is %s update in  %s namespace", "serivce", instance.Name, instance.Namespace))
			}
		} else {
			if err := r.Delete(ctx, s); err != nil {
				return nil, err
			}
			r.Recorder.Event(instance, corev1.EventTypeNormal, fmt.Sprintf("%s-controller", "service"), fmt.Sprintf("type is %s name is %s delete in  %s namespace", "serivce", instance.Name, instance.Namespace))
		}
	}
	return &workloadsv1alpha1.ServiceStatus{
		ServiceIP: s.Spec.ClusterIP,
	}, nil
}

// 处理wk status的 func
func (r *WorkloadReconciler) workloadStatus(instance *workloadsv1alpha1.Workload, dgStatus *workloadsv1alpha1.DeploymentGroupStatus, svcStatus *workloadsv1alpha1.ServiceStatus, ctx context.Context) error {
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
	instance.Status = s
	// todo
	return r.Status().Update(ctx, instance)
}

// 处理wk phase的 func
func (r *WorkloadReconciler) workloadPhase(ctx context.Context, instance *workloadsv1alpha1.Workload, phase workloadsv1alpha1.Phase) error {
	instance.Status.Phase = phase
	return r.Status().Update(ctx, instance)
}
