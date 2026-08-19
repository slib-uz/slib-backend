package asynqtask

import (
	"github.com/hibiken/asynq"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/async"
	"slib.uz/src/infrastructure/asynqtask/mapper"
	"slib.uz/src/infrastructure/asynqtask/middleware"
)

type AsyncServerImpl struct {
	server *asynq.Server
	mux    *asynq.ServeMux

	// middlewares
	taskResultMiddleware *middleware.TaskResultMiddleware
}

// @inject
func NewAsyncServerImpl(server *asynq.Server, mux *asynq.ServeMux, resultMiddleware *middleware.TaskResultMiddleware) async.AsyncServer {
	return &AsyncServerImpl{server: server, mux: mux, taskResultMiddleware: resultMiddleware}
}

func (this *AsyncServerImpl) Init() {
	this.mux.Use(mapper.ToAsynqMiddleware(this.taskResultMiddleware.Wrap))
}

func (this *AsyncServerImpl) Run() error {
	return this.server.Run(this.mux)
}

func (this *AsyncServerImpl) Use(middlewares ...async.AsyncTaskMiddleware) {
	for _, m := range middlewares {
		this.mux.Use(mapper.ToAsynqMiddleware(m))
	}
}

func (this *AsyncServerImpl) HandlerFunc(taskName enum.Task, handler async.AsyncTaskHandler) {
	this.mux.HandleFunc(string(taskName), mapper.ToAsynqHandler(handler).ProcessTask)
}
