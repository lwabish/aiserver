package sadtalker

import (
	"errors"
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

func (s *controller) Setup() {
	panic(errors.New("not implemented"))
}
