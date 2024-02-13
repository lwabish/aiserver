package controllers

import (
	"github.com/jinzhu/gorm"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"github.com/sirupsen/logrus"
)

var (
	BaseCtl = &BaseController{}
)

type BaseController struct {
	DB *gorm.DB
	Q  *utils.TaskQueue
	L  *logrus.Logger
}
type BaseControllerCfg struct {
	DB *gorm.DB
	Q  *utils.TaskQueue
	L  *logrus.Logger
}

func Inject(cfg *BaseControllerCfg) {
	BaseCtl.Q = cfg.Q
	BaseCtl.L = cfg.L
	BaseCtl.DB = cfg.DB
}

func (b *BaseController) UpdateTaskStatus(uid string, status models.TaskStatus) {
	r := b.DB.
		Model(&models.Task{}).
		Where("uid = ?", uid).
		Update("status", status)
	if r.Error != nil || r.RowsAffected != 1 {
		b.L.Warnf("update task status error: %v", r.Error)
	}
}
