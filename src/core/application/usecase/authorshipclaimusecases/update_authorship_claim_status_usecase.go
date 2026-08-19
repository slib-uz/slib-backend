package authorshipclaimusecases

import (
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/logger"
)

type UpdateAuthorshipClaimStatusUseCase struct {
	claimRepo                  repository.AuthorshipClaimRepository
	articleRepo                repository.ArticleRepository
	publishedArticleRepository repository.PublishedArticleRepository
	authorRepo                 repository.AuthorRepository
	roiGateway                 gateway.ROIGateway
	frontendURL                string
	log                        *logger.AsyncLogger
}

// @inject
func NewUpdateAuthorshipClaimStatusUseCase(
	claimRepo repository.AuthorshipClaimRepository,
	articleRepo repository.ArticleRepository,
	publishedArticleRepository repository.PublishedArticleRepository,
	authorRepo repository.AuthorRepository,
	roiGateway gateway.ROIGateway,
	configAdapter conf.ConfigAdapter,
	log *logger.AsyncLogger,
) *UpdateAuthorshipClaimStatusUseCase {
	return &UpdateAuthorshipClaimStatusUseCase{
		claimRepo:                  claimRepo,
		articleRepo:                articleRepo,
		publishedArticleRepository: publishedArticleRepository,
		authorRepo:                 authorRepo,
		roiGateway:                 roiGateway,
		frontendURL:                configAdapter.GetFrontendURL(),
		log:                        log,
	}
}

func (this *UpdateAuthorshipClaimStatusUseCase) Execute(user *entity.UserBasicEntity, input *entity.UpdateAuthorshipClaimStatusInputEntity) error {

	if err := this.checkUserAccess(user); err != nil {
		return err
	}

	claim, err := this.claimRepo.GetByID(input.ClaimID)
	if err != nil {
		return err
	}

	if claim.Status != enum.ClaimStatusPending {
		return response.NewFailResponse(400, "claim status is not pending")
	}

	switch input.Status {
	case enum.ClaimStatusRejected:
		if input.RejectReason == "" {
			return response.NewFailResponse(400, "reject reason is empty")
		}
		claim.RejectReason = input.RejectReason

	case enum.ClaimStatusApproved:
		// Add user to co-authors

		// Find or create author
		author, err := this.authorRepo.GetByScienceID(claim.Sender.ScienceID)
		if errors.Is(err, response.NotFoundError) {
			// Create author
			newAuthor := &entity.AuthorEntity{
				ScienceID: claim.Sender.ScienceID,
				FullName:  claim.Sender.FullName,
			}
			savedAuthor, err := this.authorRepo.SaveAuthor(newAuthor)
			if err != nil {
				return err
			}
			author = savedAuthor
		} else if err != nil {
			return err
		}

		if err := this.articleRepo.AddCoAuthor(claim.ArticleID, author.ID); err != nil {
			return err
		}
	default:
		return response.NewFailResponse(400, "invalid status")
	}

	claim.Status = input.Status
	claim.ReviewedByID = &input.ReviewerID
	now := time.Now()
	claim.ReviewedAt = &now

	err = this.claimRepo.Update(claim)
	if err != nil {
		return err
	}

	if input.Status == enum.ClaimStatusApproved {
		article, err := this.publishedArticleRepository.GetByIDWithRelations(claim.ArticleID)
		if err != nil {
			return err
		}

		if err := this.sendToROI(article); err != nil {
			// Xato ataylab yutiladi: da'vo allaqachon tasdiqlangan va
			// bazaga yozilgan, ROI'ga yuborish esa yordamchi qadam.
			this.log.Error("ROI'ga maqola yuborilmadi",
				zap.Uint("article_id", article.ID),
				zap.Uint("claim_id", input.ClaimID),
				zap.Error(err))
		}
	}

	return nil
}

func (this *UpdateAuthorshipClaimStatusUseCase) checkUserAccess(user *entity.UserBasicEntity) error {

	if user.IsAdmin {
		return nil
	}

	for _, role := range user.Roles {
		if role.Role == enum.RoleAdmin || role.Role == enum.RoleChiefEditor {
			return nil
		}
	}
	return response.PermissionDeniedError
}

func (this *UpdateAuthorshipClaimStatusUseCase) sendToROI(article *entity.ArticleEntity) error {
	updateArticle := mapper.ArticleEntityToROIArticleUpdateEntity(article, fmt.Sprintf("%s/article/%d", this.frontendURL, article.ID))
	if updateArticle.ROI == "" {
		return errors.New("Invalid ROI")
	}
	return this.roiGateway.UpdateArticle(updateArticle)
}
