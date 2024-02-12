package main

import (
	"fmt"
	"github.com/lwabish/cloudnative-ai-server/controllers/sadtalker"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
)

func StartWorker(queue *utils.TaskQueue) {
	for task := range queue.Chan() {
		queue.TaskOut()
		processTask(task)
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
