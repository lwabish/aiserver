package sadtalker

import (
	"github.com/lwabish/cloudnative-ai-server/config"
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"github.com/lwabish/cloudnative-ai-server/models"
	"sync"
	"time"
)

var (
	StHdl = newHandler()
)

type handler struct {
	*handlers.BaseHandler
	workerParam  map[string]*taskParam
	workerFunc   func(task *models.Task, p *taskParam) error
	pythonPath   string
	JobNamespace string
	sync.Mutex
}

func newHandler() *handler {
	for handlers.BaseHdl == nil {
		time.Sleep(100 * time.Millisecond)
	}
	h := &handler{
		BaseHandler: handlers.BaseHdl,
		workerParam: make(map[string]*taskParam),
	}
	return h
}

func (s *handler) Setup(cfg *config.Config) {
	if s.C == nil {
		s.workerFunc = s.createJob
	} else {
		s.workerFunc = s.invokeSadTalker
		s.pythonPath = cfg.BareMetal.SadTalker.PythonPath
	}
}
