package journal_applications_usecases

import (
	"errors"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type CheckISSNUsecase struct {
	repository repository.JournalApplicationRepository
}

// @inject
func NewCheckISSNUsecase(journalApplicationRepository repository.JournalApplicationRepository) *CheckISSNUsecase {
	return &CheckISSNUsecase{repository: journalApplicationRepository}
}

func (this *CheckISSNUsecase) Execute(issn string, issnType string) error {
	app, err := this.repository.GetByISSN(issn, issnType)

	if err != nil {
		if errors.Is(err, response.NotFoundError) {
			return nil
		}
		return err
	}

	if app.Status == enum.StatusPending {
		return response.NewOptionalResponse(200, response.CodeApplicationAlreadyExists, app, "")
	}

	if app.Status == enum.StatusAccepted {
		return response.NewOptionalResponse(200, response.CodeJournalAlreadyExists, app, "")
	}

	return nil
}
