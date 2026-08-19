package mapper

import (
	"time"

	"gorm.io/gorm"
	entity "slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

// calculateRatingAvg вычисляет средний рейтинг на основе суммы и количества
func calculateRatingAvg(ratingSum, ratingCount int64) float64 {
	if ratingCount > 0 {
		return float64(ratingSum) / float64(ratingCount)
	}
	return 0.0
}

func JournalCreateEntityToModel(journal *entity.JournalCreateEntity) *models.JournalModel {

	studyFields := make([]*models.StudyFieldModel, len(journal.StudyFieldIDs))
	for i, id := range journal.StudyFieldIDs {
		studyFields[i] = &models.StudyFieldModel{Model: gorm.Model{ID: id}}
	}

	languages := make([]*models.LanguageModel, len(journal.LanguageIDs))
	for i, id := range journal.LanguageIDs {
		languages[i] = &models.LanguageModel{Model: gorm.Model{ID: id}}
	}

	// Parse date string to time.Time
	var establishedDate time.Time
	if journal.DateOfEstablished != "" {
		establishedDate, _ = time.Parse("2006-01-02", journal.DateOfEstablished)
	}

	model := models.NewJournalModel(
		ToGormJson(journal.Name),
		ToGormJson(journal.ShortName),
		journal.ISSNPaper,
		journal.ISSNOnline,
		ToGormJson(journal.Description),
		ToGormJson(journal.Rule),
		ToGormJson(journal.ArticlePublishConditions),
		establishedDate,
		journal.Website,
		journal.CertificateFile,
		journal.Email,
		ToGormJson(journal.Address),
		journal.PhoneNumber,
		journal.CoverImageFile,
		journal.PublisherID,
		studyFields,
		journal.PublishingPrice,
		journal.SellingPrice,
		journal.AccessType,
		journal.OAKCertificateFile,
		journal.PeerReviewMethod,
		journal.SubmissionAccess,
		journal.CommentAccess,
		ToGormJson(journal.SocialNetworks),
		journal.SupportLink,
	)
	model.Languages = languages
	model.RegionID = journal.RegionID
	model.DistrictID = journal.DistrictID
	return model
}

func JournalModelToEntity(it *models.JournalModel) *entity.JournalEntity {
	var publisher *entity.PublisherEntity
	var studyFields []*entity.StudyFieldEntity
	var indexes []*entity.JournalIndexingEntity
	var languages []*entity.LanguageEntity

	if len(it.Languages) > 0 {
		languages = make([]*entity.LanguageEntity, len(it.Languages))
		for i, lang := range it.Languages {
			languages[i] = LanguageModelToEntity(lang)
		}
	}

	if it.Publisher != nil {
		publisher = PublisherModelToEntity(it.Publisher)
	}

	if it.StudyFields != nil {
		studyFields = make([]*entity.StudyFieldEntity, len(it.StudyFields))
		for i, item := range it.StudyFields {
			studyFields[i] = StudyFieldModelToEntity(item)
		}
	}
	if it.Indexes != nil {
		indexes = make([]*entity.JournalIndexingEntity, len(it.Indexes))
		for i, index := range it.Indexes {
			indexes[i] = JournalIndexingModelToEntity(index)
		}
	}

	ratingAvg := calculateRatingAvg(it.RatingSum, it.RatingCount)

	journalEntity := entity.NewJournalEntity(
		it.ID,
		FromGormJson[map[string]string](it.Name),
		FromGormJson[map[string]string](it.ShortName),
		it.ISSNPaper,
		it.ISSNOnline,
		FromGormJson[map[string]string](it.Description),
		FromGormJson[map[string]string](it.ArticlePublishConditions),
		FromGormJson[map[string]string](it.Rule),
		it.EstablishedDate,
		it.Website,
		indexes,
		it.CertificateFile,
		it.Email,
		FromGormJson[map[string]any](it.Address),
		it.PhoneNumber,
		it.CoverImageFile,
		it.PublisherID,
		publisher,
		studyFields,
		it.PublishingPrice,
		it.SellingPrice,
		it.AccessType,
		it.OAKCertificateFile,
		languages,
		it.PeerReviewMethod,
		it.IsActive,
		it.IsAccepted,
		it.RatingCount,
		it.RatingSum,
		ratingAvg,
		it.ViewsCount,
		0,
		it.SubmissionAccess,
		it.CommentAccess,
		FromGormJson[map[string]any](it.SocialNetworks),
		it.SupportLink,
	)
	journalEntity.RegionID, journalEntity.DistrictID, journalEntity.Region, journalEntity.District = MapRegionDistrictFromModel(
		it.RegionID,
		it.DistrictID,
		it.Region,
		it.District,
	)
	return journalEntity
}

func JournalModelToEntityWithApplicationID(it *models.JournalModel, applicationID uint) *entity.JournalEntity {
	var publisher *entity.PublisherEntity
	var studyFields []*entity.StudyFieldEntity
	var indexes []*entity.JournalIndexingEntity
	var languages []*entity.LanguageEntity

	if len(it.Languages) > 0 {
		languages = make([]*entity.LanguageEntity, len(it.Languages))
		for i, lang := range it.Languages {
			languages[i] = LanguageModelToEntity(lang)
		}
	}

	if it.Publisher != nil {
		publisher = PublisherModelToEntity(it.Publisher)
	}

	if len(it.StudyFields) > 0 {
		studyFields = make([]*entity.StudyFieldEntity, len(it.StudyFields))
		for i, item := range it.StudyFields {
			studyFields[i] = StudyFieldModelToEntity(item)
		}
	}
	if len(it.Indexes) > 0 {
		indexes = make([]*entity.JournalIndexingEntity, len(it.Indexes))
		for i, index := range it.Indexes {
			indexes[i] = JournalIndexingModelToEntity(index)
		}
	}

	ratingAvg := calculateRatingAvg(it.RatingSum, it.RatingCount)

	journalEntity := entity.NewJournalEntity(
		it.ID,
		FromGormJson[map[string]string](it.Name),
		FromGormJson[map[string]string](it.ShortName),
		it.ISSNPaper,
		it.ISSNOnline,
		FromGormJson[map[string]string](it.Description),
		FromGormJson[map[string]string](it.ArticlePublishConditions),
		FromGormJson[map[string]string](it.Rule),
		it.EstablishedDate,
		it.Website,
		indexes,
		it.CertificateFile,
		it.Email,
		FromGormJson[map[string]any](it.Address),
		it.PhoneNumber,
		it.CoverImageFile,
		it.PublisherID,
		publisher,
		studyFields,
		it.PublishingPrice,
		it.SellingPrice,
		it.AccessType,
		it.OAKCertificateFile,
		languages,
		it.PeerReviewMethod,
		it.IsActive,
		it.IsAccepted,
		it.RatingCount,
		it.RatingSum,
		ratingAvg,
		it.ViewsCount,
		applicationID,
		it.SubmissionAccess,
		it.CommentAccess,
		FromGormJson[map[string]any](it.SocialNetworks),
		it.SupportLink,
	)
	journalEntity.RegionID, journalEntity.DistrictID, journalEntity.Region, journalEntity.District = MapRegionDistrictFromModel(
		it.RegionID,
		it.DistrictID,
		it.Region,
		it.District,
	)
	return journalEntity
}

func JournalModelToBasicEntity(journal *models.JournalModel) *entity.JournalBasicEntity {
	var publisher *entity.PublisherEntity
	if journal.Publisher != nil {
		publisher = PublisherModelToEntity(journal.Publisher)
	}

	var studyFields []*entity.StudyFieldEntity
	if journal.StudyFields != nil {
		studyFields = make([]*entity.StudyFieldEntity, len(journal.StudyFields))
		for i, item := range journal.StudyFields {
			studyFields[i] = StudyFieldModelToEntity(item)
		}
	}

	var indexes []*entity.JournalIndexingEntity
	if journal.Indexes != nil {
		indexes = make([]*entity.JournalIndexingEntity, len(journal.Indexes))
		for i, index := range journal.Indexes {
			indexes[i] = JournalIndexingModelToEntity(index)
		}
	}

	ratingAvg := calculateRatingAvg(journal.RatingSum, journal.RatingCount)

	journalEntity := entity.NewJournalBasicEntity(
		journal.ID,
		FromGormJson[map[string]string](journal.Name),
		FromGormJson[map[string]string](journal.ShortName),
		FromGormJson[map[string]string](journal.Description),
		journal.ISSNPaper,
		journal.ISSNOnline,
		journal.EstablishedDate,
		journal.Website,
		FromGormJson[map[string]any](journal.Address),
		journal.PhoneNumber,
		journal.CoverImageFile,
		journal.PublisherID,
		publisher,
		studyFields,
		journal.AccessType,
		journal.IsActive,
		journal.RatingCount,
		journal.RatingSum,
		ratingAvg,
		journal.ViewsCount,
		journal.OAKCertificateFile,
		indexes,
	)
	journalEntity.RegionID, journalEntity.DistrictID, journalEntity.Region, journalEntity.District = MapRegionDistrictFromModel(
		journal.RegionID,
		journal.DistrictID,
		journal.Region,
		journal.District,
	)
	return journalEntity
}
