package main

import (
	"flag"

	container "slib.uz/cmd/container"
)

var mode string

func init() {
	flag.StringVar(&mode, "mode", "http", "run mode: http | task-worker | scheduler")
}

func main() {
	flag.Parse()
	switch mode {
	case "http":
		startHttpServer()
	case "task-worker":
		startJobServer()
	case "scheduler":
		startSchedulerServer()
	}

}

func startHttpServer() {
	app := container.InitHttpServer()
	app.Init()
	app.Start()
}

func startJobServer() {
	app := container.InitAsynqServer()
	app.Init()
	app.Start()
}

func startSchedulerServer() {
	app := container.InitScheduler()
	app.Init()
	app.Start()
}
