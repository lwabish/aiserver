package main

import (
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"log"
	"time"
)

func StartWorker(queue *utils.TaskQueue) {
	for task := range queue.Chan() {
		queue.TaskOut()
		processTask(task)
	}
}

func processTask(task *models.Task) {
	// 这里添加任务处理逻辑
	log.Printf("processing task: %+v\n", task)
	time.Sleep(5 * time.Second)
}
