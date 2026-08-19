package journal_applications_usecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type JournalCreateUseCase struct {
	repository repository.JournalApplicationRepository
	language   repository.LanguageRepository
	storage    storage.FileStorage
	checkIssn  *CheckISSNUsecase
}

// @inject
func NewJournalCreateUseCase(repository repository.JournalApplicationRepository, language repository.LanguageRepository, storage storage.FileStorage, checkIssn *CheckISSNUsecase) *JournalCreateUseCase {
	return &JournalCreateUseCase{repository: repository, language: language, storage: storage, checkIssn: checkIssn}
}

func (this *JournalCreateUseCase) Execute(request *entity.JournalCreateEntity) error {

	if this.checkIssn == nil {
		return this.repository.Create(request)
	}

	if request.ISSNPaper != nil {
		if err := this.checkIssn.Execute(*request.ISSNPaper, "issn_paper"); err != nil {
			return err
		}
	}
	if request.ISSNOnline != nil {
		if err := this.checkIssn.Execute(*request.ISSNOnline, "issn_online"); err != nil {
			return err
		}
	}

	return this.repository.Create(request)
}
