package usecase

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type SpellCheckUsecase struct {
	applicationRepository repository.ApplicationRepository
	task                  *tasks.SpellCheckTaskSender
	resultRepository      repository.SpellCheckResultRepository
	billingGateway        gateway.SlibBillingGateway
}

// @inject
func NewSpellCheckUsecase(
	applicationRepository repository.ApplicationRepository,
	task *tasks.SpellCheckTaskSender,
	resultRepository repository.SpellCheckResultRepository,
	billingGateway gateway.SlibBillingGateway,
) *SpellCheckUsecase {
	return &SpellCheckUsecase{
		applicationRepository: applicationRepository,
		task:                  task,
		resultRepository:      resultRepository,
		billingGateway:        billingGateway,
	}
}

func (this *SpellCheckUsecase) Execute(userID, reviewStageID, applicationID uint) (uint, error) {

	app, err := this.applicationRepository.GetWithArticleAndJournal(applicationID)
	if err != nil {
		return 0, err
	}

	// Check journal account balance from billing
	balance, err := this.billingGateway.GetJournalServiceBalance(app.JournalID, enum.ServiceSavodxon)
	if err != nil {
		return 0, err
	}

	if !balance.IsAvailable() {
		return 0, response.InsufficientBalance
	}

	contentFile := app.Article.ContentFile
	spellcheckResult, err := this.resultRepository.Create(entity.NewSpellCheckResultEntity(0, reviewStageID, nil, app.ID, nil, app.JournalID, nil, contentFile, nil, 0, userID, nil, nil))

	if err != nil {
		return 0, err
	}

	if err := this.task.Run(spellcheckResult); err != nil {
		return 0, err
	}
	return spellcheckResult.ID, nil

}
