package handlers

import (
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"sync"
)

var (
	BaseHdl = newController()
)

type TaskParam interface {
	String() string
}

type WorkerFunc func(task *models.Task, tp TaskParam) error

type BaseHandler struct {
	DB *gorm.DB
	Q  *utils.TaskQueue
	L  *logrus.Logger
	// nil -> bare metal
	// non nil -> k8s
	C *kubernetes.Clientset

	// fixme: 老化gc
	// some svc -> map[uid]TaskParam
	TaskParams  map[string]map[string]TaskParam
	TaskWorkers map[string]WorkerFunc
	sync.Mutex
}

func newController() *BaseHandler {
	b := &BaseHandler{
		TaskParams:  map[string]map[string]TaskParam{},
		TaskWorkers: map[string]WorkerFunc{},
	}
	return b
}

type BaseHandlerCfg struct {
	DB *gorm.DB
	Q  *utils.TaskQueue
	L  *logrus.Logger
	C  *kubernetes.Clientset
}

func (b *BaseHandler) Setup(cfg *BaseHandlerCfg) {
	b.DB = cfg.DB
	b.Q = cfg.Q
	b.L = cfg.L
}

func (b *BaseHandler) SetupCloudNative(cfg *BaseHandlerCfg) {
	b.C = cfg.C
}

func (b *BaseHandler) SetWorkerFunc(t string, f WorkerFunc) {
	b.Lock()
	defer b.Unlock()
	b.TaskWorkers[t] = f
}

func (b *BaseHandler) SetTaskParam(t string, uid string, param TaskParam) {
	b.Lock()
	defer b.Unlock()
	if b.TaskParams[t] == nil {
		b.TaskParams[t] = make(map[string]TaskParam)
	}
	b.TaskParams[t][uid] = param
}

func (b *BaseHandler) GetTaskParam(t, uid string) TaskParam {
	b.Lock()
	defer b.Unlock()
	if b.TaskParams[t] == nil {
		return nil
	}
	return b.TaskParams[t][uid]
}

func (b *BaseHandler) UpdateTaskStatus(uid string, status models.TaskStatus) {
	r := b.DB.
		Model(&models.Task{}).
		Where("uid = ?", uid).
		Update("status", status)
	if r.Error != nil || r.RowsAffected != 1 {
		b.L.Warnf("update task status error: %v", r.Error)
	}
}

func (b *BaseHandler) SaveTaskResult(uid string, result string) {
	r := b.DB.
		Model(&models.Task{}).
		Where("uid = ?", uid).
		Update("result", result)
	if r.Error != nil || r.RowsAffected != 1 {
		b.L.Warnf("update task result error: %v", r.Error)
	}
}

func (b *BaseHandler) Process(task *models.Task) {
	b.UpdateTaskStatus(task.Uid, models.TaskStatusRunning)
	var err error
	defer func() {
		if err != nil {
			b.L.Errorf("Process task failed: %s %s", task.Uid, err.Error())
			b.UpdateTaskStatus(task.Uid, models.TaskStatusFailed)
		}
	}()
	p := b.GetTaskParam(task.Type, task.Uid)
	b.L.Infof("Processing task: %s %s", task.Uid, p)
	workerFunc := b.TaskWorkers[task.Type]
	err = workerFunc(task, p)
}
