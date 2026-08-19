package userusecases

import (
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type ProfileUpdateUseCase struct {
	repository repository.UserProfileRepository
	storage    storage.FileStorage
}

func NewProfileUpdateUseCase(repository repository.UserProfileRepository, storage storage.FileStorage) *ProfileUpdateUseCase {
	return &ProfileUpdateUseCase{repository: repository, storage: storage}
}

func (this *ProfileUpdateUseCase) Execute(id uint, bio *string, email string, _photo *string) error {
	//var photo *string
	//if _photo != nil {
	//	uploadedPhoto, err := this.uploadFile(*_photo)
	//	if err != nil {
	//		return fmt.Errorf("failed to upload photo: %w", err)
	//	}
	//	photo = &uploadedPhoto
	//}
	return this.repository.Update(id, bio, email, _photo)
}

//func (this *ProfileUpdateUseCase) uploadFile(photo string) (string, error) {
//	objectPath, err := this.storage.UploadIfExists(enum.FolderImages, photo, enum.BucketPublic)
//	if err != nil {
//		return "", fmt.Errorf("failed to upload image: %w", err)
//	}
//
//	return objectPath, nil
//}
