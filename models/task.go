package models

import (
	"github.com/jinzhu/gorm"
)

// Task 模型定义，对应之前Python示例中的数据库表结构
type Task struct {
	gorm.Model
	Uid    string     `json:"uid"`
	Type   string     `json:"type"`
	Result string     `json:"result"`
	Status TaskStatus `json:"status"`
}
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusSuccess   TaskStatus = "success"
	TASK_STATUS_MISSING TaskStatus = "missing_result"
	TaskStatusFailed    TaskStatus = "failed"
)
