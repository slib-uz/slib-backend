package tasks

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
)

type SpellCheckTaskSender struct {
	publisher.TaskPublisher
}

// @inject
func NewSpellCheckTaskSender(taskPublisher publisher.TaskPublisher) *SpellCheckTaskSender {
	return &SpellCheckTaskSender{TaskPublisher: taskPublisher}
}

func (this *SpellCheckTaskSender) Run(spellcheckResult *entity2.SpellCheckResultEntity) error {
	data := entity2.NewTaskEntity[any](enum.TaskSpellcheck, spellcheckResult)
	return this.Publish(data, 0)
}
