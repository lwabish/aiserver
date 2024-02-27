package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"path"
)

var (
	BaseHdl = &BaseHandler{}
)

type BaseHandler struct {
	DB *gorm.DB
	Q  *utils.TaskQueue
	L  *logrus.Logger
	// nil -> bare metal
	// non nil -> k8s
	C *kubernetes.Clientset
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

// DownloadResult 下载任务结果
func (b *BaseHandler) DownloadResult(c *gin.Context) {
	fileName := c.PostForm("filename")
	c.FileAttachment(path.Join(utils.ResultDir, fileName), fileName)
}

// GetTaskStatus 查询任务状态
func (b *BaseHandler) GetTaskStatus(c *gin.Context) {
	var task models.Task
	if err := b.DB.Where("uid = ?", c.PostForm("task_id")).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     task.Uid,
		"status": task.Status,
		"index":  b.Q.FindTaskPosition(task.Uid),
		"result": task.Result,
	})
}
