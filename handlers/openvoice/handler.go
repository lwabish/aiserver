package openvoice

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
		"mp3", "wav", "flac", "aac", "m4a", "ogg", "opus", "wma",
	}
)

// UploadFile 上传文件并创建任务
func (h *handler) UploadFile(c *gin.Context) {
	text := c.PostForm("text")
	audio, err := c.FormFile("audio")
	if audio.Filename == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source/target file not found"})
		return
	}

	if !utils.IsAllowedExtension(audio.Filename, allowedUploadExtensions) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file extension not allowed"})
		return
	}

	audioPath := filepath.Join(utils.UploadDir, audio.Filename)
	if err = c.SaveUploadedFile(audio, audioPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	p := &taskParam{
		text:      text,
		audioPath: audioPath,
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
