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
	workerParam map[string]*taskParam
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

// InjectCfg 如果sub controller有配置，通过main包调用注入配置和其他依赖
func (s *controller) InjectCfg() {

}
