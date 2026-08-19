package authusecases

import (
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type UserLoginUseCase struct {
	gateway    gateway.ScienceIdOAuthGateway
	auth       *service.UserAuthTokenService
	repository repository.UserRepository
	//doctorateUpdateTask      *tasks.DoctorateUpdateTask
	//degreeUpdateTask         *tasks.DegreeUpdateTask
	//jobsUpdateTask           *tasks.JobsUpdateTask
	//academicDegreeUpdateTask *tasks.AcademicDegreeUpdateTask
	//academicTitleUpdateTask  *tasks.AcademicTitleUpdateTask
	//photoMigrationTask       *tasks.UserPhotoMigrationTask
}

// NewUserLoginUseCase is a `CONSTRUCTOR` for UserLoginUseCase
// @inject
func NewUserLoginUseCase(
	gateway gateway.ScienceIdOAuthGateway,
	auth *service.UserAuthTokenService,
	repository repository.UserRepository,
	// doctorateUpdateTask *tasks.DoctorateUpdateTask,
	// degreeUpdateTask *tasks.DegreeUpdateTask,
	// jobsUpdateTask *tasks.JobsUpdateTask,
	// academicDegreeUpdateTask *tasks.AcademicDegreeUpdateTask,
	// academicTitleUpdateTask *tasks.AcademicTitleUpdateTask,
	// photoMigrationTask *tasks.UserPhotoMigrationTask,
) *UserLoginUseCase {
	return &UserLoginUseCase{
		gateway:    gateway,
		auth:       auth,
		repository: repository,
		//doctorateUpdateTask:      doctorateUpdateTask,
		//degreeUpdateTask:         degreeUpdateTask,
		//jobsUpdateTask:           jobsUpdateTask,
		//academicDegreeUpdateTask: academicDegreeUpdateTask,
		//academicTitleUpdateTask:  academicTitleUpdateTask,
		//photoMigrationTask:       photoMigrationTask,
	}
}

func (this *UserLoginUseCase) Execute(code string) (*entity.AuthTokenEntity, error) {

	user, err := this.gateway.GetUserByCode(code)
	if err != nil {
		return nil, err
	}

	scienceID := []rune(user.ScienceID)

	cityCode := string(scienceID[1:3])
	user.City = enum.CityCode(cityCode)

	existingUser, err := this.repository.GetByScienceId(user.ScienceID)
	isNewUser := err != nil || existingUser == nil
	hasAcademicDegree := existingUser != nil && string(existingUser.AcademicDegree) != ""

	if isNewUser || !hasAcademicDegree {
		academicDegreeCode := string(scienceID[0:1])
		user.AcademicDegree = enum.AcademicDegreeCode(academicDegreeCode)
	} else {
		user.AcademicDegree = existingUser.AcademicDegree
	}

	//hasExistingPhoto := existingUser != nil && existingUser.Photo != ""
	//scienceIDPhoto := user.Photo
	//if hasExistingPhoto {
	//user.Photo = existingUser.Photo
	//}

	user, err = this.repository.Save(user)
	if err != nil {
		return nil, err
	}

	//shouldMigratePhoto := (isNewUser || !hasExistingPhoto) && scienceIDPhoto != ""
	//if shouldMigratePhoto {
	//	_ = this.photoMigrationTask.Run(user.ID, scienceIDPhoto)
	//}

	return this.auth.GenerateToken(user.ID), nil

}
