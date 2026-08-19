package repository

import "slib.uz/src/core/domain/entity"

type TagRepository interface {
	GetOrCreateTags(tags entity.TagNamesByLang) ([]uint, error)
}
