package tasks

import (
	"go.uber.org/zap"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
	"slib.uz/src/infrastructure/logger"
)

type SetRoiSenderTask struct {
	publisher publisher.TaskPublisher
	log       *logger.AsyncLogger
}

// @inject
func NewSetRoiSenderTask(publisher publisher.TaskPublisher, log *logger.AsyncLogger) *SetRoiSenderTask {
	return &SetRoiSenderTask{publisher: publisher, log: log}
}

type SetRoiSenderPayload struct {
	User          *entity.UserBasicEntity
	ApplicationID uint
	File          string
}

// Run arizani ROI o'rnatish navbatiga qo'yadi. Navbat nosozligi
// chaqiruvchiga qaytarilmaydi — ROI o'rnatish yordamchi qadam va u
// asosiy oqimni yiqitmasligi kerak. Shuning uchun Run har doim
// nil qaytaradi; imzodagi error kelajakdagi kengaytirish uchun qolgan.
func (this *SetRoiSenderTask) Run(payload SetRoiSenderPayload) error {
	task := entity.NewTaskEntity[any](enum.TaskSetRoi, payload)

	if err := this.publisher.Publish(task, 5); err != nil {
		this.log.Error("ROI o'rnatish vazifasi navbatga qo'yilmadi",
			zap.Uint("application_id", payload.ApplicationID),
			zap.Error(err))
	}

	return nil
}
