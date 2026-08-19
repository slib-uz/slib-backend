package logusecases

import (
	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/tasks"
)

type AdminAlertUseCase struct {
	sendTelegramAlertTask *tasks.SendTelegramAlertTask
}

// @inject
func NewAdminAlertUseCase(sendTelegramAlertTask *tasks.SendTelegramAlertTask) *AdminAlertUseCase {
	return &AdminAlertUseCase{sendTelegramAlertTask: sendTelegramAlertTask}
}

func (this *AdminAlertUseCase) Execute(message string) {
	if err := this.sendTelegramAlertTask.Run(message); err != nil {
		log.Error("Failed to enqueue telegram alert task: ", err)
	}
}
