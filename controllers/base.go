package controllers

import (
	"github.com/jinzhu/gorm"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"github.com/sirupsen/logrus"
)

var (
	BaseCtl = NewBaseController(&BaseControllerCfg{})
)

type BaseController struct {
	DB *gorm.DB
	Q  *utils.TaskQueue
	L  *logrus.Logger
}
type BaseControllerCfg struct {
}

func (b *BaseController) UpdateTaskStatus(task *models.Task) {
	//b.DB.Update()
}

func NewBaseController(_ *BaseControllerCfg) *BaseController {
	return &BaseController{}
}
