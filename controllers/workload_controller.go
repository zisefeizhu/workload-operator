/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	"github.com/zisefeizhu/workload-operator/controllers/template"
	"github.com/zisefeizhu/workload-operator/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

// WorkloadReconciler reconciles a Workload object
type WorkloadReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Logger   logr.Logger
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=workloads.zise.feizhu,resources=workloads,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=workloads.zise.feizhu,resources=workloads/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=workloads.zise.feizhu,resources=workloads/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Workload object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile

/*
  1、svc 先于 工作负载创建
  2、工作负载适配:deployment";"statefulSet";"daemonSet";"job";"cronJob
  3、 目前尚未对任务重新入队列进行处理
*/

func (r *WorkloadReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//logger := log.FromContext(ctx)
	_ = r.Logger.WithValues("workloads", req.NamespacedName)
	instance := &workloadsv1alpha1.Workload{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// svc 处理逻辑
	svcStatus, err := r.svc(instance, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	// 工作负载 处理逻辑
	dgStatus, err := r.deploymentGroup(instance, ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}

	// wk 的 status 处理
	err = r.workLoadStatus(instance, dgStatus, svcStatus, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 2 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WorkloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&workloadsv1alpha1.Workload{}).
		Complete(r)
}

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
	return utils.WorkloadStatusProcessor(w), nil
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
func (r *WorkloadReconciler) workLoadStatus(instance *workloadsv1alpha1.Workload, dgStatus *workloadsv1alpha1.DeploymentGroupStatus, svcStatus *workloadsv1alpha1.ServiceStatus, ctx context.Context) error {
	// status
	s := workloadsv1alpha1.WorkloadStatus{}
	s.DeploymentGroupStatus = *dgStatus
	s.ServiceStatus = *svcStatus
	instance.Status = s
	// todo
	err := r.Status().Update(ctx, instance)
	if err != nil {
		return err
	}
	return nil
}
