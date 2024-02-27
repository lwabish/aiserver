package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"net/http"
	"path"
)

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
