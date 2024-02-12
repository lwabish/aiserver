package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

// NewLogger 创建并配置新的日志实例
func NewLogger(level logrus.Level) *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{}
	log.Level = level
	log.Out = os.Stdout
	return log
}
