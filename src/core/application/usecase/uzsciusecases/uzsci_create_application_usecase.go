package uzsciusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/utils"
)

type UzSciCreateApplicationUseCase struct {
	gateway     gateway.UzsciGateway
	journalRepo repository.JournalRepository
}

// @inject
func NewUzSciCreateApplicationUseCase(gateway gateway.UzsciGateway, journalRepo repository.JournalRepository) *UzSciCreateApplicationUseCase {
	return &UzSciCreateApplicationUseCase{
		gateway:     gateway,
		journalRepo: journalRepo,
	}
}

func (this *UzSciCreateApplicationUseCase) Execute(periodID uint, journalID uint, answers []entity.UzSciApplicationAnswerEntity) error {
	journal, err := this.journalRepo.FindByID(journalID)
	if err != nil {
		return err
	}

	issnPrint, issnOnline := extractValidISSNs(journal)
	if issnPrint == "" && issnOnline == "" {
		return response.NewFailResponse(400, "journal has no valid ISSN")
	}

	externalJournal, err := this.gateway.GetJournalByISSN(issnPrint, issnOnline)
	if err != nil {
		return err
	}
	if externalJournal == nil {
		return response.NewFailResponse(404, "uzsci journal not found")
	}

	return this.gateway.CreateApplication(periodID, externalJournal.ID, answers)
}

func extractValidISSNs(journal *entity.JournalEntity) (issnPrint, issnOnline string) {
	if journal.ISSNPaper != nil && utils.IsValidISSN(*journal.ISSNPaper) {
		issnPrint = *journal.ISSNPaper
	}
	if journal.ISSNOnline != nil && utils.IsValidISSN(*journal.ISSNOnline) {
		issnOnline = *journal.ISSNOnline
	}
	return issnPrint, issnOnline
}
