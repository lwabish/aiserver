package sadtalker

import (
	"context"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *handler) createJob(task *models.Task, _ *taskParam) error {
	j := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      task.Uid,
			Namespace: s.JobNamespace,
			Annotations: map[string]string{
				utils.TaskTypeKey: TaskType,
			},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
					Volumes: []v1.Volume{
						{Name: "data",
							VolumeSource: v1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: "ai-server",
								}}}},
					Containers: []v1.Container{{
						Name:            "sad-talker",
						Image:           "ccr.ccs.tencentyun.com/lwabish/sadtalker",
						ImagePullPolicy: v1.PullIfNotPresent,
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								"nvidia.com/gpu": *resource.NewQuantity(1, resource.DecimalSI),
							},
						},
						VolumeMounts: []v1.VolumeMount{{
							Name:      "data",
							MountPath: "/app/SadTalker/results",
							SubPath:   "sad-talker/result",
						}},
						//Args: []string{},//todo
					}},
				},
			},
		},
	}
	_, err := s.C.BatchV1().Jobs(s.JobNamespace).Create(context.TODO(), j, metav1.CreateOptions{})
	return err
}
