package k8s

import (
	"k8s.io/api/core/v1"
)

func toPods(pods []v1.Pod) []Pod {
	var out []Pod
	for _, pod := range pods {
		out = append(out, toPod(pod))
	}
	return out
}

type Pod struct {
	Name      string
	Namespace string
}

func toPod(pod v1.Pod) Pod {
	return Pod{
		Name:      pod.Name,
		Namespace: pod.Namespace,
	}
}
