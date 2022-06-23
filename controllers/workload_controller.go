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
	"github.com/zisefeizhu/workload-operator/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	workloadsv1alpha1 "github.com/zisefeizhu/workload-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
*/

func (r *WorkloadReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//logger := log.FromContext(ctx)
	_ = r.Logger.WithValues("workloads", req.NamespacedName)
	workload := &workloadsv1alpha1.Workload{}
	if err := r.Get(ctx, req.NamespacedName, workload); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	err := r.svc(workload, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	// 获取工作负载模版
	app := NewWorkload(workload).Template()
	if err := controllerutil.SetControllerReference(workload, app.(metav1.Object), r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	d := NewWorkload(workload).Found()
	if err := r.Get(ctx, req.NamespacedName, d.(client.Object)); err != nil {
		if errors.IsNotFound(err) {
			err := r.Create(ctx, app.(client.Object))
			if err != nil {
				r.Logger.Error(err, "create app failed")
				return ctrl.Result{RequeueAfter: 2 * time.Second}, err
			}
			r.Recorder.Event(workload, corev1.EventTypeNormal, fmt.Sprintf("%s-controller", workload.Spec.Type), fmt.Sprintf("type is %s name is %s create in  %s namespace", workload.Spec.Type, workload.Name, workload.Namespace))
		}
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
	} else {
		// 这里update 需要优化
		// todo
		// 需要比对 新旧工作负载的spec 字段 ，类此
		//if !reflect.DeepEqual(deploy.Spec, found.Spec) {
		//	found.Spec = deploy.Spec
		//}
		if err := r.Update(ctx, app.(client.Object)); err != nil {
			r.Logger.Error(err, "update app failed")
			return ctrl.Result{}, err
		}
		r.Recorder.Event(workload, corev1.EventTypeNormal, fmt.Sprintf("%s-controller", workload.Spec.Type), fmt.Sprintf("type is %s name is %s  update in %s namespace", workload.Spec.Type, workload.Name, workload.Namespace))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WorkloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&workloadsv1alpha1.Workload{}).
		Complete(r)
}

func (r *WorkloadReconciler) svc(workload *workloadsv1alpha1.Workload, ctx context.Context) error {
	// service的处理
	service := utils.NewService(workload)
	err := controllerutil.SetControllerReference(workload, service, r.Scheme)
	if err != nil {
		return err
	}
	s := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{Name: workload.Name, Namespace: workload.Namespace}, s); err != nil {
		if errors.IsNotFound(err) && workload.Spec.EnableService {
			if err := r.Create(ctx, service); err != nil {
				r.Logger.Error(err, "create service failed")
				return err
			}
			r.Recorder.Event(workload, corev1.EventTypeNormal, fmt.Sprintf("%s-controller", "service"), fmt.Sprintf("type is %s name is %s create in  %s namespace", "serivce", workload.Name, workload.Namespace))
		}
		if !errors.IsNotFound(err) && workload.Spec.EnableService {
			return err
		}
	} else {
		if workload.Spec.EnableService {
			currIP := s.Spec.ClusterIP
			s.Spec.ClusterIP = currIP
			service.Spec.ClusterIP = currIP
			if !reflect.DeepEqual(s.Spec, service.Spec) {
				s.Spec = service.Spec
				if err := r.Update(ctx, s); err != nil {
					r.Logger.Error(err, "update service failed")
					return err
				}
				r.Recorder.Event(workload, corev1.EventTypeNormal, fmt.Sprintf("%s-controller", "service"), fmt.Sprintf("type is %s name is %s update in  %s namespace", "serivce", workload.Name, workload.Namespace))
			}
		} else {
			if err := r.Delete(ctx, s); err != nil {
				return err
			}
			r.Recorder.Event(workload, corev1.EventTypeNormal, fmt.Sprintf("%s-controller", "service"), fmt.Sprintf("type is %s name is %s delete in  %s namespace", "serivce", workload.Name, workload.Namespace))
		}
	}
	return nil
}
