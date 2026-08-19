package journalusecases

import (
	"strings"

	"slib.uz/src/core/application/response"
	entity2 "slib.uz/src/core/domain/entity"
	enum2 "slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalListUseCase struct {
	repository repository.JournalRepository
}

// @inject
func NewJournalListUseCase(repository repository.JournalRepository) *JournalListUseCase {
	return &JournalListUseCase{repository: repository}
}

func (this JournalListUseCase) Execute(page, size int, submissionAccess int, oakRegistered *bool, publisherId uint, name, description, issn, publisherName *string, languageIds, studyFieldIds []uint, fromYear, toYear *int, indexingTypes []string, sortBy, order string) (*entity2.PagingEntity[entity2.JournalBasicEntity], error) {
	submissionAccessEnum := parseAccessType(submissionAccess)

	sortBy, order, err := ValidateSortOrder(sortBy, order)
	if err != nil {
		return nil, err
	}

	var indexingTypesEnum []enum2.IndexingType
	for _, indexingType := range indexingTypes {
		indexingTypesEnum = append(indexingTypesEnum, enum2.IndexingType(indexingType))
	}
	paging, err := this.repository.GetListByPage(page, size, submissionAccessEnum, oakRegistered, publisherId, name, description, issn, publisherName, languageIds, studyFieldIds, fromYear, toYear, indexingTypesEnum, sortBy, order)

	if err != nil {
		return nil, err
	}

	return paging, nil
}

func parseAccessType(value int) enum2.AccessType {
	switch value {
	case int(enum2.PublicAccessType):
		return enum2.PublicAccessType
	case int(enum2.PrivateAccessType):
		return enum2.PrivateAccessType
	default:
		return 0
	}
}

// ValidateSortOrder faqat yo'nalishni normallashtiradi va tekshiradi.
// Maydon nomining ro'yxatga mosligini repozitoriydagi JournalSortFields
// hal qiladi — ustun nomlari persistence qatlamining tushunchasi, va
// ikkita ro'yxat saqlash ularning bir-biridan uzoqlashishiga olib keladi.
func ValidateSortOrder(sortBy, order string) (string, string, error) {
	if sortBy == "" {
		if order != "" {
			return "", "", response.NewFailResponse(400, "invalid order: sort_by is required")
		}
		return "", "", nil
	}

	if order == "" {
		return sortBy, "desc", nil
	}

	order = strings.ToLower(order)
	if order != "asc" && order != "desc" {
		return "", "", response.NewFailResponse(400, "invalid order: must be asc or desc")
	}

	return sortBy, order, nil
}
