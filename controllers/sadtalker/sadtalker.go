package sadtalker

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lwabish/cloudnative-ai-server/controllers"
	"github.com/lwabish/cloudnative-ai-server/models"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	StCtl                   = newController()
	allowedUploadExtensions = []string{
		".png", ".jpg", ".jpeg", ".gif", ".mp3", ".wav", ".m4a", ".mp4",
	}
)

const (
	uploadDir = "uploads"
	resultDir = "results"
)

type controller struct {
	*controllers.BaseController
	workerParam map[string]*taskParam
	sync.Mutex
}

func newController() *controller {
	for controllers.BaseCtl == nil {
		time.Sleep(100 * time.Millisecond)
	}
	return &controller{
		BaseController: controllers.BaseCtl,
		workerParam:    make(map[string]*taskParam),
	}
}

// InjectCfg 如果sub controller有配置，通过main包调用注入配置和其他依赖
func (s *controller) InjectCfg() {

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
func (s *controller) GetTaskStatus(c *gin.Context) {
	var task models.Task
	if err := s.DB.Where("uid = ?", c.PostForm("task_id")).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     task.Uid,
		"status": task.Status,
		"index":  s.Q.FindTaskPosition(task.Uid),
	})
}

// DownloadResult 下载任务结果
func (s *controller) DownloadResult(c *gin.Context) {
	c.FileAttachment(resultDir, c.PostForm("filename"))
}

func isAllowedExtension(fileName string) bool {
	for _, extension := range allowedUploadExtensions {
		if strings.ToLower(filepath.Ext(fileName)) == extension {
			return true
		}
	}
	return false
}
