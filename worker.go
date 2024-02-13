package main

import (
	"fmt"
	"github.com/lwabish/cloudnative-ai-server/controllers/sadtalker"
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
	switch task.Type {
	case sadtalker.TaskType:
		sadtalker.StCtl.Process(task)
	default:
		panic(fmt.Errorf("unknown task type: %s", task.Type))
	}
}
