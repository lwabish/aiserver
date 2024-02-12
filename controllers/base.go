package controllers

import (
	"github.com/jinzhu/gorm"
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

func NewBaseController(_ *BaseControllerCfg) *BaseController {
	return &BaseController{}
}
