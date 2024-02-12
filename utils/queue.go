package utils

import (
	"github.com/lwabish/cloudnative-ai-server/models"
	"sync"
)

type TaskQueue struct {
	tasksChan chan *models.Task
	tasksList []*models.Task // 使用slice维护任务列表
	lock      sync.Mutex
}

func NewTaskQueue(bufferSize int) *TaskQueue {
	return &TaskQueue{
		tasksChan: make(chan *models.Task, bufferSize),
		tasksList: make([]*models.Task, 0),
	}
}

func (q *TaskQueue) Enqueue(task *models.Task) {
	q.lock.Lock()
	q.tasksList = append(q.tasksList, task) // 将任务添加到列表末尾
	q.lock.Unlock()

	q.tasksChan <- task // 将任务发送到通道，以便工作线程处理
}

func (q *TaskQueue) TaskOut() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.tasksList = q.tasksList[:len(q.tasksList)]
}

func (q *TaskQueue) FindTaskPosition(taskID uint) int {
	q.lock.Lock()
	defer q.lock.Unlock()

	for i, task := range q.tasksList {
		if task.ID == taskID {
			return i // 找到任务，返回其在tasksList中的位置
		}
		// todo: 设置上限，如果超了返回-2，前端显示>1h
	}
	return -1 // 未找到任务
}

func (q *TaskQueue) Chan() chan *models.Task {
	return q.tasksChan
}
