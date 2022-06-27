package tool

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func TimeToKubernetes() string {
	t := metav1.Time{Time: time.Now().UTC()}
	return t.Format("2006-01-02T15:04:05Z")
}
