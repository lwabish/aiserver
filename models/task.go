package models

import (
	"github.com/jinzhu/gorm"
)

// Task 模型定义，对应之前Python示例中的数据库表结构
type Task struct {
	gorm.Model
	Result string `json:"result"`
	Status string `json:"status"`
}
