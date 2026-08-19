package wire

import (
	"github.com/google/wire"
	"slib.uz/src/entrypoint/asynctask"
	app "slib.uz/src/entrypoint/presentation/app"
	"slib.uz/src/entrypoint/scheduler"
)

func InitHttpServer() *app.App {
	wire.Build(ProviderSet)
	return nil
}

func InitAsynqServer() *asynctask.App {
	wire.Build(ProviderSet)
	return nil
}

func InitScheduler() *scheduler.App {
	wire.Build(ProviderSet)
	return nil
}
