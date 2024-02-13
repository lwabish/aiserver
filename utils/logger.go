package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

// NewLogger 创建并配置新的日志实例
func NewLogger(level logrus.Level) *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// 这里返回的第一个字符串是函数名，我们这里忽略
			// 第二个字符串是文件名和行号
			return "", fmt.Sprintf(" [%s:%d]:", path.Base(f.File), f.Line)
		},
	}
	log.SetReportCaller(true)
	log.Level = level
	log.Out = os.Stdout
	return log
}
