package repository

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
	infraError "slib.uz/src/infrastructure/errors"
)

type ArticleFileRepositoryImpl struct {
	*BaseRepository
	articleRepository repository.ArticleRepository
	storage           storage.FileStorage
}

// @inject
func NewArticleFileRepository(baseRepository *BaseRepository, articleRepository repository.ArticleRepository, storage storage.FileStorage) repository.ArticleFileRepository {
	return &ArticleFileRepositoryImpl{BaseRepository: baseRepository, articleRepository: articleRepository, storage: storage}
}

func (this *ArticleFileRepositoryImpl) GetFile(filePath string) (*entity.ArticleFileEntity, error) {

	expires := 7 * 24 * time.Hour
	fileUrl, fErr := this.storage.PresignedURL(enum.BucketPrivate, filePath, expires)
	if fErr != nil {
		return nil, infraError.Wrap(fErr)
	}
	binaryFile, bErr := this.storage.GetObject(enum.BucketPrivate, filePath)
	if bErr != nil {
		return nil, infraError.Wrap(bErr)
	}

	return entity.NewArticleFileEntity(fileUrl, binaryFile), nil

}
