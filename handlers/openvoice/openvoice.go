package openvoice

import (
	"github.com/lwabish/cloudnative-ai-server/config"
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"time"
)

var (
	Handler = newHandler()
)

type handler struct {
	*handlers.BaseHandler

	extraArgs []string

	// bare metal
	projectPath string
	pythonPath  string
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

func (h *handler) Setup(cfg *config.Config) {
	if h.C == nil {
		h.SetWorkerFunc(TaskType, h.invoke)
		c := cfg.BareMetal.OpenVoice
		h.pythonPath = c.PythonPath
		h.projectPath = c.ProjectPath
		h.extraArgs = c.ExtraArgs
	}
}
