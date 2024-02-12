package models

import (
	"github.com/jinzhu/gorm"
)

// Task 模型定义，对应之前Python示例中的数据库表结构
type Task struct {
	gorm.Model
	Uid    string `json:"uid"`
	Result string `json:"result"`
	Status string `json:"status"`
}

const (
	TASK_STATUS_PENDING = "pending"
	TASK_STATUS_RUNNING = "running"
	TASK_STATUS_SUCCESS = "success"
	TASK_STATUS_MISSING = "missing_result"
	TASK_STATUS_FAILED  = "failed"
)
