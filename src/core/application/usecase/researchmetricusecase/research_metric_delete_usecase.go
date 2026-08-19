package researchmetricusecase

import (
	"errors"

	"gorm.io/gorm"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/ports/repository"
)

type ResearchMetricDeleteUseCase struct {
	repository repository.ResearchMetricRepository
}

// @inject
func NewResearchMetricDeleteUseCase(repository repository.ResearchMetricRepository) *ResearchMetricDeleteUseCase {
	return &ResearchMetricDeleteUseCase{repository: repository}
}

func (this *ResearchMetricDeleteUseCase) Execute(userID, id uint) error {
	err := this.repository.DeleteByIDAndUserID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewFailResponse(404, "Research metric not found")
		}
		return err
	}
	return nil
}
