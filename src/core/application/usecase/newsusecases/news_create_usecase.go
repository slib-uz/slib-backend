package newsusecases

import (
	"context"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type NewsCreateUseCase struct {
	newsRepository repository.NewsRepository
	storage        storage.FileStorage
}

// @inject
func NewNewsCreateUseCase(newsRepository repository.NewsRepository, storage storage.FileStorage) *NewsCreateUseCase {
	return &NewsCreateUseCase{newsRepository: newsRepository, storage: storage}
}

func (this *NewsCreateUseCase) Execute(ctx context.Context, news *entity.NewsEntity) (*entity.NewsEntity, error) {
	//if path, err := this.uploadImage(news.Image); err != nil {
	//	return nil, err
	//} else {
	//	news.Image = path
	//}

	if err := this.newsRepository.Create(news); err != nil {
		return nil, err
	}

	return news, nil
}

//func (this *NewsCreateUseCase) uploadImage(localPath string) (string, error) {
//	if localPath == "" {
//		return "", nil
//	}
//	path, err := this.storage.UploadIfExists(enum.FolderNews, localPath, enum.BucketPublic)
//	if err != nil {
//		return "", err
//	}
//	return path, nil
//}
