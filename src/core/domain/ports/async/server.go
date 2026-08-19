package async

import (
	"slib.uz/src/core/domain/entity/enum"
)

type AsyncServer interface {
	Init()
	Run() error
	Use(middlewares ...AsyncTaskMiddleware)
	HandlerFunc(taskName enum.Task, handler AsyncTaskHandler)
}
