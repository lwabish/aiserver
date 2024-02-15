package controllers

import (
	"context"
	"fmt"
	"io"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func jobSucceeded(job *batchv1.Job) bool {
	complete := false
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobComplete {
			complete = true
		}
	}
	return complete && *job.Spec.Completions == job.Status.Succeeded
}

func getPodLogString(client *kubernetes.Clientset, jobName string, namespace string) (string, error) {
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", jobName),
	})
	if err != nil {
		return "", err
	}
	var podName string
	if len(pods.Items) > 0 {
		podName = pods.Items[0].Name
	}

	logReq := client.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{})
	logs, err := logReq.Stream(context.TODO())
	if err != nil {
		return "", err
	}
	defer func() { _ = logs.Close() }()

	bytes, err := io.ReadAll(logs)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
