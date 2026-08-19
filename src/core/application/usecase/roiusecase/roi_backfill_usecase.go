package roiusecase

import (
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/domain/ports/repository"
)

// ROIBackfillUseCase — ROI olinmay qolib ketgan published maqolalar uchun
// publish task'ni navbatga qo'yadi (qayta-qayta xavfsiz chaqirsa bo'ladi).
type ROIBackfillUseCase struct {
	articleRepository    repository.ArticleRepository
	publishRoiSenderTask *tasks.PublishRoiSenderTask
}

// @inject
func NewROIBackfillUseCase(
	articleRepository repository.ArticleRepository,
	publishRoiSenderTask *tasks.PublishRoiSenderTask,
) *ROIBackfillUseCase {
	return &ROIBackfillUseCase{
		articleRepository:    articleRepository,
		publishRoiSenderTask: publishRoiSenderTask,
	}
}

type ROIBackfillResult struct {
	Enqueued int    `json:"enqueued"`
	IDs      []uint `json:"ids"`
}

func (this *ROIBackfillUseCase) Execute() (*ROIBackfillResult, error) {
	ids, err := this.articleRepository.PublishedIDsWithoutROI()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		_ = this.publishRoiSenderTask.Run(tasks.PublishRoiSenderPayload{ArticleID: id})
	}

	return &ROIBackfillResult{
		Enqueued: len(ids),
		IDs:      ids,
	}, nil
}
