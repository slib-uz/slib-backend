package researchmetricusecase

import (
	"errors"
	"fmt"
	"net/url"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type ResearchMetricCreateOrUpdateUseCase struct {
	repository repository.ResearchMetricRepository
}

// @inject
func NewResearchMetricCreateOrUpdateUseCase(repository repository.ResearchMetricRepository) *ResearchMetricCreateOrUpdateUseCase {
	return &ResearchMetricCreateOrUpdateUseCase{repository: repository}
}

func (this *ResearchMetricCreateOrUpdateUseCase) Execute(userID uint, entity *entity.ResearchMetricEntity) error {
	if err := this.isValidURL(entity.ProfileUrl); err != nil {
		return err
	}

	if err := this.validate(entity.Source); err != nil {
		return err
	}

	entity.UserID = userID
	err := this.repository.UpdateOrCreate(userID, entity.Source, entity)

	if err != nil {
		return err
	}

	return nil
}

func (this *ResearchMetricCreateOrUpdateUseCase) isValidURL(urlString string) error {
	parsed, err := url.Parse(urlString)
	if err != nil {
		return errors.New("invalid URL format")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("URL must start with http or https")
	}

	if parsed.Host == "" {
		return errors.New("URL must have a valid host")
	}

	return nil
}

func (this *ResearchMetricCreateOrUpdateUseCase) validate(source enum.ResearchMetricEnum) error {
	if source != enum.ScopusIndexing && source != enum.WosIndexing {
		return fmt.Errorf("invalid source: must be either %s or %s", enum.ScopusIndexing, enum.WosIndexing)
	}
	return nil
}
