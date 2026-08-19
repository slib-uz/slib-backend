package reportuscases

import (
	"strings"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type ReportListUseCase struct {
	repository repository.ReportRepository
}

// @inject
func NewReportListUseCase(repository repository.ReportRepository) *ReportListUseCase {
	return &ReportListUseCase{repository: repository}
}

func (this ReportListUseCase) Execute(page, pageSize int, ordering string, targetTypeStr string) (*entity2.PagingEntity[entity2.ReportEntity], error) {
	var targetType enum.ReportType
	if targetTypeStr != "" {
		switch strings.ToLower(targetTypeStr) {
		case string(enum.ArticleReport):
			targetType = enum.ArticleReport
		case string(enum.JournalReport):
			targetType = enum.JournalReport
		default:
			targetType = ""
		}
	}
	paging, err := this.repository.GetByPaging(page, pageSize, ordering, targetType)
	if err != nil {
		return nil, err
	}
	return paging, nil
}
