package publisherusecases

import (
	"errors"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type PublisherCreateUseCase struct {
	repository repository.PublisherRepository
	storage    storage.FileStorage
}

// @inject
func NewPublisherCreateUseCase(repository repository.PublisherRepository, storage storage.FileStorage) *PublisherCreateUseCase {
	return &PublisherCreateUseCase{repository: repository, storage: storage}
}

func (this *PublisherCreateUseCase) Execute(publisher *entity.PublisherEntity) error {
	//if path, err := this.uploadLogo(publisher.Logo); err != nil {
	//	return err
	//} else {
	//	publisher.Logo = path
	//}

	if _, err := this.repository.GetByTin(publisher.Tin); err == nil {
		return response.NewOptionalResponse(400, response.CodePublisherTinAlreadyExists, nil, "publisher with this TIN already exists")
	} else if !errors.Is(err, response.NotFoundError) {
		return err
	}

	if err := this.repository.Create(publisher); err != nil {
		return err
	}
	return nil
}
