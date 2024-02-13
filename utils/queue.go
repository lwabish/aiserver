package utils

import (
	"github.com/gammazero/deque"
	"github.com/lwabish/cloudnative-ai-server/models"
)

type TaskQueue struct {
	*deque.Deque[*models.Task]
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		deque.New[*models.Task](),
	}
}

func (q *TaskQueue) Enqueue(task *models.Task) {
	q.PushBack(task)
}

func (q *TaskQueue) FindTaskPosition(uid string) int {
	return q.Index(func(task *models.Task) bool {
		return task.Uid == uid
	})
}
