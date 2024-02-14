package sadtalker

import (
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"sync"
	"time"
)

var (
	StHdl = newHandler()
)

type handler struct {
	*handlers.BaseHandler
	workerParam  map[string]*taskParam
	JobNamespace string
	sync.Mutex
}

func newHandler() *handler {
	for handlers.BaseHdl == nil {
		time.Sleep(100 * time.Millisecond)
	}
	return &handler{
		BaseHandler: handlers.BaseHdl,
		workerParam: make(map[string]*taskParam),
	}
}

type Cfg struct {
	JobNamespace string
}

func (s *handler) Setup(cfg *Cfg) {
	s.JobNamespace = cfg.JobNamespace
}
