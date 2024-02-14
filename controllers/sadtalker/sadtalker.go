package sadtalker

import (
	"github.com/lwabish/cloudnative-ai-server/controllers"
	"sync"
	"time"
)

var (
	StCtl = newController()
)

type controller struct {
	*controllers.BaseController
	workerParam  map[string]*taskParam
	JobNamespace string
	sync.Mutex
}

func newController() *controller {
	for controllers.BaseCtl == nil {
		time.Sleep(100 * time.Millisecond)
	}
	return &controller{
		BaseController: controllers.BaseCtl,
		workerParam:    make(map[string]*taskParam),
	}
}

type Cfg struct {
	JobNamespace string
}

func (s *controller) Setup(cfg *Cfg) {
	s.JobNamespace = cfg.JobNamespace
}
