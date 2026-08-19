package main

import (
	container "slib.uz/cmd/container"
)

func main() {
	app := container.InitAsynqServer()
	app.Init()
	app.Start()
}
