package sadtalker

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lwabish/cloudnative-ai-server/controllers"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"net/http"
)

var (
	stCtl = newSadTalkerController(&cfg{})
)

func newSadTalkerController(_ *cfg) *controller {
	return &controller{
		controllers.BaseCtl,
	}
}

type controller struct {
	*controllers.BaseController
}

type cfg struct {
}

// UploadFile 上传文件并创建任务
func (s *controller) UploadFile(c *gin.Context) {
	// 模拟文件保存逻辑，实际应用中需要保存上传的文件
	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file part"})
		return
	}
	// 假设文件保存成功，创建任务
	task := models.Task{
		Status: "pending",
	}
	//s.Create(&task)
	c.JSON(http.StatusCreated, gin.H{"task_id": task.ID})
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
