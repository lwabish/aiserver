package sadtalker

import (
	"github.com/lwabish/cloudnative-ai-server/models"
	"time"
)

type taskParam struct {
	photo string
	audio string
}

func (s *controller) Process(task *models.Task) {
	s.UpdateTaskStatus(task.Uid, models.TaskStatusRunning)
	var err error
	defer func() {
		if err != nil {
			s.UpdateTaskStatus(task.Uid, models.TaskStatusFailed)
		}
	}()
	p := s.getParam(task.Uid)
	s.L.Infof("Processing sad talker task %v+ %v+", task, p)

	//todo
	time.Sleep(10 * time.Second)

	s.UpdateTaskStatus(task.Uid, models.TaskStatusSuccess)
}

func (s *controller) getParam(uid string) *taskParam {
	s.Lock()
	defer s.Unlock()
	return s.workerParam[uid]
}

func (s *controller) setParam(uid string, param *taskParam) {
	s.Lock()
	defer s.Unlock()
	s.workerParam[uid] = param
}
