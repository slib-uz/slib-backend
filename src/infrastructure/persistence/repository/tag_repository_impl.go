package repository

import (
	"errors"

	"gorm.io/gorm"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/utils"
	"slib.uz/src/infrastructure/persistence/models"
)

type TagRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewTagRepository(baseRepository *BaseRepository) repository.TagRepository {
	return &TagRepositoryImpl{BaseRepository: baseRepository}
}

func (this *TagRepositoryImpl) GetOrCreateTags(tags entity.TagNamesByLang) ([]uint, error) {
	return getOrCreateTagsByLang(this.db(), tags)
}

func getOrCreateTagsByLang(db *gorm.DB, tags entity.TagNamesByLang) ([]uint, error) {
	pairs, err := utils.NormalizeTagNamesByLang(tags)
	if err != nil {
		return nil, response.NewFailResponse(400, err.Error())
	}
	ids := make([]uint, 0, len(pairs))
	seen := make(map[uint]bool, len(pairs))
	for _, pair := range pairs {
		id, err := findOrCreateTag(db, pair.Name, pair.Lang)
		if err != nil {
			return nil, err
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}
	return ids, nil
}

func findOrCreateTag(db *gorm.DB, name, lang string) (uint, error) {
	var tag models.TagModel
	err := db.Where("lang = ? AND lower(name) = lower(?)", lang, name).First(&tag).Error
	if err == nil {
		return tag.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	tag = models.TagModel{Name: name, Lang: lang}
	if err := db.Create(&tag).Error; err != nil {
		return 0, err
	}
	return tag.ID, nil
}
