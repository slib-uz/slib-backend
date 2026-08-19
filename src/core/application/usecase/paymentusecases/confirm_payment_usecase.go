package paymentusecases

import (
	"time"

	entity2 "slib.uz/src/core/domain/entity"
	enum2 "slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/utils"
	"slib.uz/src/infrastructure/config"
)

type ConfirmPaymentUseCase struct {
	articlePurchaseRepository repository.ArticlePurchaseRepository
	reviewStageRepository     repository.ReviewStageRepository
	env                       *config.Config
}

// @inject
func NewConfirmPaymentUseCase(articlePurchaseRepository repository.ArticlePurchaseRepository, reviewStageRepository repository.ReviewStageRepository, env *config.Config) *ConfirmPaymentUseCase {
	return &ConfirmPaymentUseCase{articlePurchaseRepository: articlePurchaseRepository, reviewStageRepository: reviewStageRepository, env: env}
}

func (this *ConfirmPaymentUseCase) Execute(applicationID uint) error {

	paymentStage, err := this.reviewStageRepository.GetByStage(applicationID, enum2.PaymentStage)
	if err != nil {
		return err
	}

	if paymentStage.Status == enum2.StatusAccepted {
		return nil
	}

	now := time.Now()

	paymentStage.Status = enum2.StatusAccepted
	paymentStage.ReviewedAt = &now

	nextStage := entity2.NewReviewStageEntity(
		0,
		paymentStage.ApplicationID,
		nil,
		enum2.PublishStage,
		enum2.StatusPending,
		nil,
		nil,
		nil,
		nil,
		now,
		nil,
		nil,
		false,
		nil,
		nil,
		nil,
		utils.MakeDeadline(time.Now(), this.env.ReviewDeadlineDays),
		nil,
	)

	return this.reviewStageRepository.UpdateToNext(paymentStage, nextStage)
}
