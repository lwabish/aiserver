package roop

import (
	"github.com/lwabish/cloudnative-ai-server/config"
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"github.com/lwabish/cloudnative-ai-server/models"
	"sync"
	"time"
)

var (
	Handler = newHandler()
)

type handler struct {
	*handlers.BaseHandler
	workerParam map[string]*taskParam
	workerFunc  func(task *models.Task, p *taskParam) error
	sync.Mutex

	extraArgs []string

	// bare metal
	projectPath string
	pythonPath  string
}

func (h *handler) Process(task *models.Task) {
	h.UpdateTaskStatus(task.Uid, models.TaskStatusRunning)
	var err error
	defer func() {
		if err != nil {
			h.L.Errorf("Process roop task failed: %s %s", task.Uid, err.Error())
			h.UpdateTaskStatus(task.Uid, models.TaskStatusFailed)
		}
	}()
	p := h.getParam(task.Uid)
	h.L.Infof("Processing roop task: %s %s", task.Uid, p)
	err = h.workerFunc(task, p)
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

func (h *handler) Setup(cfg *config.Config) {
	if h.C == nil {
		h.workerFunc = h.invoke
		c := cfg.BareMetal.Roop
		h.pythonPath = c.PythonPath
		h.projectPath = c.ProjectPath
		h.extraArgs = c.ExtraArgs
	}
}
