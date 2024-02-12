package sadtalker

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lwabish/cloudnative-ai-server/controllers"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

var (
	StCtl                   = newSadTalkerController(&cfg{})
	allowedUploadExtensions = []string{
		"png", "jpg", "jpeg", "gif", "mp3", "wav", "m4a", "mp4",
	}
)

const (
	uploadDir = ""
)

func newSadTalkerController(_ *cfg) *controller {
	return &controller{
		BaseController: controllers.BaseCtl,
		workerParam:    make(map[string]taskParam),
	}
}

type controller struct {
	*controllers.BaseController
	workerParam map[string]taskParam
	sync.Mutex
}

type cfg struct {
}

// UploadFile 上传文件并创建任务
func (s *controller) UploadFile(c *gin.Context) {
	photo, pErr := c.FormFile("photo")
	audio, aErr := c.FormFile("audio")
	if pErr != nil || aErr != nil || photo.Filename == "" || audio.Filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo/audio file not found"})
		return
	}

	if !isAllowedExtension(photo.Filename) || !isAllowedExtension(audio.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file extension not allowed"})
		return
	}

	photoPath := filepath.Join(uploadDir, photo.Filename)
	audioPath := filepath.Join(uploadDir, audio.Filename)
	if pErr, aErr = c.SaveUploadedFile(photo, photoPath), c.SaveUploadedFile(audio, audioPath); pErr != nil || aErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	task := &models.Task{
		Uid:    "",
		Type:   TaskType,
		Status: models.TaskStatusPending,
	}
	s.DB.Create(task)
	s.Q.Enqueue(task)
	c.JSON(http.StatusCreated, gin.H{"task_id": task.Uid})
}

// GetTaskStatus 查询任务状态，适应POST方法
func GetTaskStatus(c *gin.Context, db *gorm.DB, taskQueue *utils.TaskQueue) {
	taskID := c.PostForm("task_id")
	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	position := taskQueue.FindTaskPosition(task.ID)

	c.JSON(http.StatusOK, gin.H{
		"id":       task.ID,
		"status":   task.Status,
		"position": position,
	})
}

// DownloadResult 下载任务结果，适应POST方法
func DownloadResult(c *gin.Context, db *gorm.DB) {
	taskID := c.PostForm("task_id")
	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	if task.Status != "success" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task is not completed"})
		return
	}
	// 实际应用中应该根据task.Result提供文件下载
	c.JSON(http.StatusOK, gin.H{"message": "Download API to be implemented"})
}

func isAllowedExtension(fileName string) bool {
	for _, extension := range allowedUploadExtensions {
		if strings.ToLower(filepath.Ext(fileName)) == extension {
			return true
		}
	}
	return false
}
