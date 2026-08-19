package main

import (
	container "slib.uz/cmd/container"
)

func main() {
	app := container.InitScheduler()
	app.Init()
	app.Start()
}
