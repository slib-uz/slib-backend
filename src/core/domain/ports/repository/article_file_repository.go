package repository

import (
	"slib.uz/src/core/domain/entity"
)

type ArticleFileRepository interface {
	GetFile(filePath string) (*entity.ArticleFileEntity, error)
}
