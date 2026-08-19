package newsusecases

import (
	"context"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type NewsUpdateUseCase struct {
	newsRepository repository.NewsRepository
	storage        storage.FileStorage
}

// @inject
func NewNewsUpdateUseCase(newsRepository repository.NewsRepository, storage storage.FileStorage) *NewsUpdateUseCase {
	return &NewsUpdateUseCase{newsRepository: newsRepository, storage: storage}
}

func (this *NewsUpdateUseCase) Execute(ctx context.Context, id uint, news *entity.NewsEntity) error {
	//existingNews, err := this.newsRepository.GetByID(id)
	//if err != nil {
	//	return err
	//}

	// Upload new image if it's different from existing one
	//if news.Image != "" && news.Image != existingNews.Image {
	//	if path, err := this.uploadImage(news.Image); err != nil {
	//		return err
	//	} else {
	//		news.Image = path
	//	}
	//}

	if err := this.newsRepository.Update(id, news); err != nil {
		return err
	}

	return nil
}

//func (this *NewsUpdateUseCase) uploadImage(localPath string) (string, error) {
//	if localPath == "" {
//		return "", nil
//	}
//	path, err := this.storage.UploadIfExists(enum.FolderNews, localPath, enum.BucketPublic)
//	if err != nil {
//		return "", err
//	}
//	return path, nil
//}
