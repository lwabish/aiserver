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

const (
	JobType = "sad-talker"
)

type taskParam struct {
	photo string
	audio string
}

func (p taskParam) String() string {
	return p.photo + "|" + p.audio
}

func (s *handler) Process(task *models.Task) {
	s.UpdateTaskStatus(task.Uid, models.TaskStatusRunning)
	var err error
	defer func() {
		if err != nil {
			s.UpdateTaskStatus(task.Uid, models.TaskStatusFailed)
		}
	}()
	p := s.getParam(task.Uid)
	s.L.Infof("Processing sad talker task:%s %s", task.Uid, p)
	err = s.createJob(task, p)
}

func (s *handler) createJob(task *models.Task, p *taskParam) error {
	j := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      task.Uid,
			Namespace: s.JobNamespace,
			Annotations: map[string]string{
				utils.TaskTypeKey: JobType,
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

func (s *handler) getParam(uid string) *taskParam {
	s.Lock()
	defer s.Unlock()
	return s.workerParam[uid]
}

func (s *handler) setParam(uid string, param *taskParam) {
	s.Lock()
	defer s.Unlock()
	s.workerParam[uid] = param
}
