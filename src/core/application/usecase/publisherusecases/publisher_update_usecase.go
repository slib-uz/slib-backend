package publisherusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type PublisherUpdateUseCase struct {
	repository repository.PublisherRepository
	storage    storage.FileStorage
}

// @inject
func NewPublisherUpdateUseCase(repository repository.PublisherRepository, storage storage.FileStorage) *PublisherUpdateUseCase {
	return &PublisherUpdateUseCase{repository: repository, storage: storage}
}

func (this PublisherUpdateUseCase) Execute(id uint, data *entity.PublisherEntity) error {
	//publisher, err := this.repository.GetByID(id)
	//if err != nil {
	//	return err
	//}

	//if data.Logo != nil && publisher.Logo == nil {
	//	objectPath, err := this.storage.UploadIfExists(enum.FolderImages, *data.Logo, enum.BucketPublic)
	//	if err != nil {
	//		return err
	//	}
	//	data.Logo = &objectPath
	//
	//} else if data.Logo != nil && publisher.Logo != nil && *data.Logo != *publisher.Logo {
	//	objectPath, err := this.storage.UploadIfExists(enum.FolderImages, *data.Logo, enum.BucketPublic)
	//	if err != nil {
	//		return err
	//	}
	//	data.Logo = &objectPath
	//}

	if err := this.repository.Update(id, data); err != nil {
		return err
	}
	return nil
}
