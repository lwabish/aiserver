package sadtalker

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"net/http"
	"path"
	"path/filepath"
)

var (
	allowedUploadExtensions = []string{
		".png", ".jpg", ".jpeg", ".gif", ".mp3", ".wav", ".m4a", ".mp4",
	}
)

// UploadFile 上传文件并创建任务
func (s *handler) UploadFile(c *gin.Context) {
	photo, pErr := c.FormFile("photo")
	audio, aErr := c.FormFile("audio")
	if pErr != nil || aErr != nil || photo.Filename == "" || audio.Filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo/audio file not found"})
		return
	}

	if !utils.IsAllowedExtension(photo.Filename, allowedUploadExtensions) || !utils.IsAllowedExtension(audio.Filename, allowedUploadExtensions) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file extension not allowed"})
		return
	}

	photoPath := filepath.Join(utils.UploadDir, photo.Filename)
	audioPath := filepath.Join(utils.UploadDir, audio.Filename)
	if pErr, aErr = c.SaveUploadedFile(photo, photoPath), c.SaveUploadedFile(audio, audioPath); pErr != nil || aErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	p := &taskParam{
		photo: photoPath,
		audio: audioPath,
	}

	task := &models.Task{
		Uid:    uuid.New().String(),
		Type:   TaskType,
		Status: models.TaskStatusPending,
	}
	s.setParam(task.Uid, p)
	if r := s.DB.Create(task); r.Error != nil {
		s.L.Warnf("create task error: %v", r.Error)
	}

	s.Q.Enqueue(task)
	c.JSON(http.StatusCreated, gin.H{"task_id": task.Uid})
}

// GetTaskStatus 查询任务状态
func (s *handler) GetTaskStatus(c *gin.Context) {
	var task models.Task
	if err := s.DB.Where("uid = ?", c.PostForm("task_id")).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     task.Uid,
		"status": task.Status,
		"index":  s.Q.FindTaskPosition(task.Uid),
		"result": task.Result,
	})
}

// DownloadResult 下载任务结果
func (s *handler) DownloadResult(c *gin.Context) {
	fileName := c.PostForm("filename")
	c.FileAttachment(path.Join(utils.ResultDir, fileName), fileName)
}
