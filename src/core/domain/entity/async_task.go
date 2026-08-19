package entity

import (
	"slib.uz/src/core/domain/entity/enum"
)

type AsyncTask struct {
	TaskID   string
	TaskType enum.Task
	Payload  []byte
}

func NewAsyncTask(taskType enum.Task, payload []byte) *AsyncTask {
	return &AsyncTask{TaskType: taskType, Payload: payload}
}
