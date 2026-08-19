package authorshipclaimusecases

import (
	"fmt"
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type CreateAuthorshipClaimUseCase struct {
	claimRepo   repository.AuthorshipClaimRepository
	articleRepo repository.ArticleRepository
	userRepo    repository.UserRepository
}

// @inject
func NewCreateAuthorshipClaimUseCase(
	claimRepo repository.AuthorshipClaimRepository,
	articleRepo repository.ArticleRepository,
	userRepo repository.UserRepository,
) *CreateAuthorshipClaimUseCase {
	return &CreateAuthorshipClaimUseCase{
		claimRepo:   claimRepo,
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

func (this *CreateAuthorshipClaimUseCase) Execute(user *entity.UserBasicEntity, input *entity.CreateAuthorshipClaimInputEntity) error {
	if len(input.ArticleIDs) == 0 {
		return response.NewFailResponse(400, "no article IDs provided")
	}

	// 2. Fetch all articles with authors
	articles, err := this.articleRepo.FindByIDsWithAuthors(input.ArticleIDs)
	if err != nil {
		return err
	}

	// Create a map for faster lookup
	articlesMap := make(map[uint]*entity.ArticleEntity)
	for _, a := range articles {
		articlesMap[a.ID] = a
	}

	// Check if all articles exist
	for _, id := range input.ArticleIDs {
		if _, ok := articlesMap[id]; !ok {
			return response.NewFailResponse(400, fmt.Sprintf("article with id %d not found", id))
		}
	}

	// 3. Fetch existing pending claims
	pendingClaims, err := this.claimRepo.FindPendingByArticleIDsAndUser(input.ArticleIDs, input.UserID)
	if err != nil {
		return err
	}
	pendingClaimsMap := make(map[uint]bool)
	for _, c := range pendingClaims {
		pendingClaimsMap[c.ArticleID] = true
	}

	// 4. Prepare claims
	claimsToCreate := make([]*entity.AuthorshipClaimEntity, 0, len(input.ArticleIDs))

	for _, articleID := range input.ArticleIDs {
		// Validations
		// a. Check if user is already an author
		article := articlesMap[articleID]
		isAuthor := false
		for _, author := range article.CoAuthors {
			if author.ScienceID != "" && author.ScienceID == user.ScienceID {
				isAuthor = true
				break
			}
		}
		if isAuthor {
			return response.NewFailResponse(400, fmt.Sprintf("author id %s is already author", user.ScienceID))
		}

		// b. Check if pending claim exists
		if pendingClaimsMap[articleID] {
			return response.NewFailResponse(400, fmt.Sprintf("pending claim for article id %d already exists", articleID))
		}

		claim := entity.NewAuthorshipClaimEntity(
			0,
			input.UserID,
			nil,
			articleID,
			nil,
			input.Comment,
			enum.ClaimStatusPending,
			time.Time{},
			time.Time{},
		)
		claimsToCreate = append(claimsToCreate, claim)
	}

	return this.claimRepo.CreateBatch(claimsToCreate)
}
