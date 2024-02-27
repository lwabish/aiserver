package sadtalker

import (
	"github.com/lwabish/cloudnative-ai-server/config"
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"time"
)

var (
	StHdl = newHandler()
)

type handler struct {
	*handlers.BaseHandler

	extraArgs []string

	// bare metal
	projectPath string
	pythonPath  string

	// cloud native
	JobNamespace string
}

func newHandler() *handler {
	for handlers.BaseHdl == nil {
		time.Sleep(100 * time.Millisecond)
	}
	h := &handler{
		BaseHandler: handlers.BaseHdl,
	}
	return h
}

func (s *handler) Setup(cfg *config.Config) {
	if s.C == nil {
		s.SetWorkerFunc(TaskType, s.invoke)
		c := cfg.BareMetal.SadTalker
		s.pythonPath = c.PythonPath
		s.projectPath = c.ProjectPath
		s.extraArgs = c.ExtraArgs
	} else {
		s.SetWorkerFunc(TaskType, s.createJob)
	}
}
