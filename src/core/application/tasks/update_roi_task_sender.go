package tasks

import (
	"go.uber.org/zap"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
	"slib.uz/src/infrastructure/logger"
)

type UpdateRoiSenderTask struct {
	publisher publisher.TaskPublisher
	log       *logger.AsyncLogger
}

// @inject
func NewUpdateRoiSenderTask(publisher publisher.TaskPublisher, log *logger.AsyncLogger) *UpdateRoiSenderTask {
	return &UpdateRoiSenderTask{publisher: publisher, log: log}
}

type UpdateRoiSenderPayload struct {
	ArticleID uint
}

// Run maqolani ROI yangilash navbatiga qo'yadi. Navbat nosozligi
// chaqiruvchiga qaytarilmaydi — ROI yangilash yordamchi qadam va u
// maqolani saqlashni yiqitmasligi kerak. Shuning uchun Run har doim
// nil qaytaradi; imzodagi error kelajakdagi kengaytirish uchun qolgan.
func (this *UpdateRoiSenderTask) Run(payload UpdateRoiSenderPayload) error {
	task := entity.NewTaskEntity[any](enum.TaskUpdateRoi, payload)

	if err := this.publisher.Publish(task, 5); err != nil {
		this.log.Error("ROI yangilash vazifasi navbatga qo'yilmadi",
			zap.Uint("article_id", payload.ArticleID),
			zap.Error(err))
	}

	return nil
}
