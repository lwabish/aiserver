package main

import (
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"time"
)

func StartWorker(queue *utils.TaskQueue) {
	for {
		time.Sleep(1 * time.Second)
		if queue.Len() != 0 {
			t := queue.PopFront()
			processTask(t)
		}
	}
}

func processTask(task *models.Task) {
	handlers.BaseHdl.Process(task)
}
