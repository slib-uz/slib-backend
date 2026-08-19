package reportuscases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ReportCreateUseCase struct {
	repository repository.ReportRepository
}

// @inject
func NewReportCreateUseCase(repository repository.ReportRepository) *ReportCreateUseCase {
	return &ReportCreateUseCase{repository: repository}
}

func (this ReportCreateUseCase) Execute(userID uint, report *entity.ReportEntity) error {
	report.ReporterID = userID
	_, err := this.repository.LoadTarget(report)
	if err != nil {
		return err
	}
	return this.repository.CreateReport(report)
}
