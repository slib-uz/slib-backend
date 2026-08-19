package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type ReportRepository interface {
	CreateReport(report *entity2.ReportEntity) error
	LoadTarget(report *entity2.ReportEntity) (interface{}, error)
	GetByPaging(page, pageSize int, ordering string, targetType enum.ReportType) (*entity2.PagingEntity[entity2.ReportEntity], error)
}
