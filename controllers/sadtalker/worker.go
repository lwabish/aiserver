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
	p := s.getParam(task.Uid)
	s.L.Infof("Processing sad talker task %v+ %v+", task, p)
	//todo
	time.Sleep(5 * time.Second)
	s.UpdateTaskStatus(task)
}

func (s *controller) getParam(uid string) taskParam {
	s.Lock()
	defer s.Unlock()
	return s.workerParam[uid]
}

func (s *controller) setParam(uid string, param taskParam) {
	s.Lock()
	defer s.Unlock()
	s.workerParam[uid] = param
}
