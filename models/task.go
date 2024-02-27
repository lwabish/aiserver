package models

import (
	"gorm.io/gorm"
)

// Task  任务模型
type Task struct {
	gorm.Model
	Uid    string     `json:"uid"`
	Type   string     `json:"type"`
	Result string     `json:"result"`
	Status TaskStatus `json:"status"`
}
type TaskStatus string

const (
	TaskStatusPending TaskStatus = "pending"
	TaskStatusRunning TaskStatus = "running"
	TaskStatusSuccess TaskStatus = "success"

	TaskStatusResultMissing TaskStatus = "result_missing"
	TaskStatusFailed        TaskStatus = "failed"
)
