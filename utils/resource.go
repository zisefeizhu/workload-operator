package utils

import (
	"github.com/zisefeizhu/workload-operator/api/v1alpha1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//func parseTemplate(templateName string, workload *v1alpha1.Workload) []byte {
//	tmpl, err := template.ParseFiles("controllers/template/" + templateName + ".yaml")
//	if err != nil {
//		panic(err)
//	}
//	b := new(bytes.Buffer)
//	err = tmpl.Execute(b, workload)
//	if err != nil {
//		panic(err)
//	}
//	return b.Bytes()
//}

func NewWorkload(workload *v1alpha1.Workload) interface{} {
	switch workload.Spec.Type {
	// 以deployment 为例
	case "deployment":
		return buildDeployment(workload)
	}
	return nil
}

func buildDeployment(workload *v1alpha1.Workload) *appv1.Deployment {
	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workload.Name,
			Namespace: workload.Namespace,
			Labels:    workload.Labels,
		},
		Spec: appv1.DeploymentSpec{
			Replicas: workload.Spec.Replicas,
			Selector: workload.Spec.Selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: workload.Labels,
				},
				Spec: workload.Spec.Template.Spec,
			},
		},
	}
}

func NewService(workload *v1alpha1.Workload) *corev1.Service {
	if workload.Spec.ServiceType == "" {
		workload.Spec.ServiceType = corev1.ServiceTypeClusterIP
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workload.Name,
			Namespace: workload.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: workload.Labels,
			Type:     workload.Spec.ServiceType,
			Ports:    workload.Spec.ServicePorts,
		},
	}
}
