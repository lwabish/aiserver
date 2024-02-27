package roop

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"net/http"
	"path/filepath"
)

var (
	allowedUploadExtensions = []string{
		".png", ".jpg", ".jpeg", ".gif",
		".mp4",
	}
)

// UploadFile 上传文件并创建任务
func (h *handler) UploadFile(c *gin.Context) {
	source, sErr := c.FormFile("source")
	target, tErr := c.FormFile("target")
	if sErr != nil || tErr != nil || source.Filename == "" || target.Filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source/target file not found"})
		return
	}

	if !utils.IsAllowedExtension(source.Filename, allowedUploadExtensions) || !utils.IsAllowedExtension(target.Filename, allowedUploadExtensions) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file extension not allowed"})
		return
	}

	sourcePath := filepath.Join(utils.UploadDir, source.Filename)
	targetPath := filepath.Join(utils.UploadDir, target.Filename)
	if sErr, tErr = c.SaveUploadedFile(source, targetPath), c.SaveUploadedFile(target, targetPath); sErr != nil || tErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	p := &taskParam{
		source: sourcePath,
		target: targetPath,
	}

	task := &models.Task{
		Uid:    uuid.New().String(),
		Type:   TaskType,
		Status: models.TaskStatusPending,
	}
	h.SetTaskParam(TaskType, task.Uid, p)
	if r := h.DB.Create(task); r.Error != nil {
		h.L.Warnf("create task error: %v", r.Error)
	}

	h.Q.Enqueue(task)
	c.JSON(http.StatusCreated, gin.H{"task_id": task.Uid})
}
