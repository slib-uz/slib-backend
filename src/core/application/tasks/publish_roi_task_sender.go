package tasks

import (
	"go.uber.org/zap"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
	"slib.uz/src/infrastructure/logger"
)

type PublishRoiSenderTask struct {
	publisher publisher.TaskPublisher
	log       *logger.AsyncLogger
}

// @inject
func NewPublishRoiSenderTask(publisher publisher.TaskPublisher, log *logger.AsyncLogger) *PublishRoiSenderTask {
	return &PublishRoiSenderTask{publisher: publisher, log: log}
}

type PublishRoiSenderPayload struct {
	ArticleID uint
}

// Run maqolani ROI nashr navbatiga qo'yadi. Navbat nosozligi
// chaqiruvchiga qaytarilmaydi — ROI yangilash yordamchi qadam va u
// maqolani saqlashni yiqitmasligi kerak. Shuning uchun Run har doim
// nil qaytaradi; imzodagi error kelajakdagi kengaytirish uchun qolgan.
func (this *PublishRoiSenderTask) Run(payload PublishRoiSenderPayload) error {
	task := entity.NewTaskEntity[any](enum.TaskPublishRoi, payload)

	if err := this.publisher.Publish(task, 5); err != nil {
		this.log.Error("ROI nashr vazifasi navbatga qo'yilmadi",
			zap.Uint("article_id", payload.ArticleID),
			zap.Error(err))
	}

	return nil
}
