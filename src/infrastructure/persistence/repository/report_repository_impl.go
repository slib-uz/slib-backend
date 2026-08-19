package repository

import (
	"errors"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
	"slib.uz/src/infrastructure/persistence/sorting"
)

type ReportRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewReportRepository(baseRepository *BaseRepository) repository.ReportRepository {
	return &ReportRepositoryImpl{BaseRepository: baseRepository}
}

func (this ReportRepositoryImpl) CreateReport(report *entity2.ReportEntity) error {
	var _model = mapper.ReportEntityToModel(report)

	if err := this.db().Model(&models.ReportModel{}).Create(_model).Error; err != nil {
		return err
	}
	return nil
}

func (this *ReportRepositoryImpl) LoadTarget(report *entity2.ReportEntity) (interface{}, error) {
	switch report.TargetType {
	case enum.ArticleReport:
		var article models.ArticleModel
		if err := this.db().First(&article, report.TargetID).Error; err != nil {
			return nil, infraError.Wrap(err)
		}
		return &article, nil
	case enum.JournalReport:
		var journal models.JournalModel
		if err := this.db().First(&journal, report.TargetID).Error; err != nil {
			return nil, infraError.Wrap(err)
		}
		return &journal, nil
	default:
		return nil, errors.New("unknown target type")
	}
}

func (r *ReportRepositoryImpl) LoadTargets(reports []*models.ReportModel) error {
	if len(reports) == 0 {
		return nil
	}

	idsByType := map[enum.ReportType][]uint{}
	for _, rep := range reports {
		idsByType[rep.TargetType] = append(idsByType[rep.TargetType], rep.TargetID)
	}

	uniq := func(in []uint) []uint {
		m := make(map[uint]struct{}, len(in))
		out := make([]uint, 0, len(in))
		for _, id := range in {
			if _, ok := m[id]; !ok {
				m[id] = struct{}{}
				out = append(out, id)
			}
		}
		return out
	}

	// load articles in one query
	if ids, ok := idsByType["article"]; ok && len(ids) > 0 {
		ids = uniq(ids)
		var articles []models.ArticleModel
		if err := r.db().Where("id IN ?", ids).Find(&articles).Error; err != nil {
			return err
		}
		am := make(map[uint]*models.ArticleModel, len(articles))
		for i := range articles {
			a := articles[i]
			am[a.ID] = &a
		}
		for _, rep := range reports {
			if rep.TargetType == "article" {
				if a, found := am[rep.TargetID]; found {
					rep.Target = a
				} else {
					rep.Target = nil
				}
			}
		}
	}

	// load journals in one query
	if ids, ok := idsByType["journal"]; ok && len(ids) > 0 {
		ids = uniq(ids)
		var journals []models.JournalModel
		if err := r.db().Where("id IN ?", ids).Find(&journals).Error; err != nil {
			return err
		}
		jm := make(map[uint]*models.JournalModel, len(journals))
		for i := range journals {
			j := journals[i]
			jm[j.ID] = &j
		}
		for _, rep := range reports {
			if rep.TargetType == "journal" {
				if j, found := jm[rep.TargetID]; found {
					rep.Target = j
				} else {
					rep.Target = nil
				}
			}
		}
	}

	return nil
}

// ReportSortFields — shikoyatlar ro'yxati uchun ruxsat etilgan tartiblash
// maydonlari. ReportModel.TableName() "report" qaytaradi — birlikda.
var ReportSortFields = sorting.New("report.created_at DESC", map[string]string{
	"created_at": "report.created_at",
})

func (this ReportRepositoryImpl) GetByPaging(page, pageSize int, ordering string, targetType enum.ReportType) (*entity2.PagingEntity[entity2.ReportEntity], error) {
	orderExpr, err := ReportSortFields.Resolve(ordering)
	if err != nil {
		return nil, err
	}

	var total int64
	var _models []*models.ReportModel

	countQuery := this.db().Model(&models.ReportModel{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	query := this.db().Model(&models.ReportModel{}).Preload("Reporter")

	if targetType != "" {
		query.Where("target_type = ?", targetType)
	}

	query = query.Order(orderExpr)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&_models).Error; err != nil {
		return nil, err
	}

	if err := this.LoadTargets(_models); err != nil {
		return nil, err
	}

	return entity2.NewPagingEntity(page, pageSize, total, mapper.ReportListModelToEntity(_models)), nil
}
