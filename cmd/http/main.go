package main

import (
	container "slib.uz/cmd/container"
)

// @title Slib.uz API Documentation
// @version 1.0
// @description API for Slib.uz
// @BasePath /api/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	app := container.InitHttpServer()
	app.Init()
	app.Start()
}
