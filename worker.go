package main

import (
	"fmt"
	"github.com/lwabish/cloudnative-ai-server/handlers/roop"
	"github.com/lwabish/cloudnative-ai-server/handlers/sadtalker"
	"github.com/lwabish/cloudnative-ai-server/models"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"time"
)

func StartWorker(queue *utils.TaskQueue) {
	for {
		time.Sleep(1 * time.Second)
		if queue.Len() != 0 {
			t := queue.PopFront()
			dispatchTask(t)
		}
	}
}

type TaskProcessor interface {
	Process(*models.Task)
}

func processTask(task *models.Task, tp TaskProcessor) {
	tp.Process(task)
}

func dispatchTask(task *models.Task) {
	switch task.Type {
	case sadtalker.TaskType:
		processTask(task, sadtalker.StHdl)
	case roop.TaskType:
		processTask(task, roop.Handler)
	default:
		panic(fmt.Errorf("unknown task type: %s", task.Type))
	}
}
