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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
  0????????????????????????ns????????????svc??????workloads ??????????????????
  1???svc ?????? ??????????????????
  2?????????????????????:deployment";"statefulSet";"daemonSet";"job";"cronJob
  3???????????????????????????spec???status?????? ???????????????
*/

func (r *WorkloadReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Logger.WithValues("workloads", req.NamespacedName)
	forget := reconcile.Result{}
	requeue := ctrl.Result{
		RequeueAfter: time.Second * 2,
	}
	instance := &workloadsv1alpha1.Workload{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return forget, client.IgnoreNotFound(err)
	}

	// ?????????????????????status

	// svc ????????????
	svcStatus, err := r.svc(instance, ctx)
	if err != nil {
		er := r.workloadPhase(ctx, instance, workloadsv1alpha1.FailedPhase)
		if er != nil {
			return requeue, er
		}
		return requeue, err
	}

	// ???????????? ????????????
	dgStatus, err := r.deploymentGroup(instance, ctx, req)
	if err != nil {
		er := r.workloadPhase(ctx, instance, workloadsv1alpha1.FailedPhase)
		if er != nil {
			return ctrl.Result{}, er
		}
		return requeue, err
	}

	// wk ??? status ??????
	err, workloadStatus := r.workloadStatus(instance, dgStatus, svcStatus, ctx)
	if err != nil {
		return requeue, err
	}
	fmt.Println(workloadStatus)

	// ????????????
	//if workloadStatus.Phase != workloadsv1alpha1.RunningPhase {
	//	// todo ????????????
	//	// ??????????????????
	//	err := r.workloadCorrectionProcessor(&workloadStatus)
	//	if err != nil {
	//		return requeue, err
	//	}
	//}

	return requeue, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WorkloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&workloadsv1alpha1.Workload{}).
		Complete(r)
}
