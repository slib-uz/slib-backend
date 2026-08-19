package etaqrizusecases

import (
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type ReSubmitForPeerReviewUseCase struct {
	gateway               gateway.ETaqrizGateway
	applicationRepository repository.ApplicationRepository
	storage               storage.FileStorage
	submissionRepository  repository.PeerReviewSubmissionRepository
}

// @inject
func NewReSubmitForPeerReviewUseCase(gateway gateway.ETaqrizGateway, applicationRepository repository.ApplicationRepository, storage storage.FileStorage, submissionRepository repository.PeerReviewSubmissionRepository) *ReSubmitForPeerReviewUseCase {
	return &ReSubmitForPeerReviewUseCase{gateway: gateway, applicationRepository: applicationRepository, storage: storage, submissionRepository: submissionRepository}
}

func (this *ReSubmitForPeerReviewUseCase) Execute(userID uint, submissionID uint, filePath string) error {
	submission, err := this.submissionRepository.GetByIDWithApplication(submissionID)
	if err != nil {
		return err
	}

	if submission.Application.UserID != userID {
		return response.PermissionDeniedError
	}

	//objectPath, err := this.uploadFile(filePath)
	//if err != nil {
	//	return err
	//}

	fileBase64, err := this.storage.GetObjectAsBase64(enum.BucketPrivate, filePath)
	if err != nil {
		return err
	}

	newDeadline, err := this.gateway.Resubmit(submission.ExternalID, fileBase64)
	if err != nil {
		return err
	}

	if err := this.applicationRepository.UpdateFile(submission.ApplicationID, filePath); err != nil {
		return err
	}

	return this.submissionRepository.UpdateStatusAndDeadline(submission.ExternalID, enum.PeerReviewStatusResubmitted, newDeadline, &submission.Deadline, time.Now().UnixMilli())
}

//func (this *ReSubmitForPeerReviewUseCase) uploadFile(newFileLocalPath string) (string, error) {
//	fileURL, err := this.storage.UploadIfExists(enum.FolderApplications, newFileLocalPath, enum.BucketPrivate)
//	if err != nil {
//		return "", err
//	}
//	return fileURL, nil
//}

//func (this *ReSubmitForPeerReviewUseCase) fileToBase64(filePath string) (string, error) {
//data, err := os.ReadFile(filePath)
//if err != nil {
//	return "", err
//}
//
//data, err := this.storage.GetObjectAsBase64(enum.)
//
//mimeType := http.DetectContentType(data)
//
//encoded := base64.StdEncoding.EncodeToString(data)
//
//return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
//}
