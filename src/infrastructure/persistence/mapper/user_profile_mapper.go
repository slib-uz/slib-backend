package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	enum2 "slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/infrastructure/persistence/models"
)

func UserMeProfileModelToEntity(profile *models.UserProfileModel) *entity2.UserMeEntity {
	var roles []*entity2.UserRoleEntity

	if profile.User != nil && profile.User.Roles != nil {
		roles = make([]*entity2.UserRoleEntity, len(profile.User.Roles))
		for i, item := range profile.User.Roles {
			roles[i] = UserRoleModelToEntity(item)
		}
	}

	var photo *string
	if profile.User != nil && profile.User.Photo != "" {
		photo = &profile.User.Photo
	} else if profile.Photo != nil {
		photo = profile.Photo
	}

	return entity2.NewUserMeEntity(
		profile.UserID,
		profile.User.FullName,
		profile.User.ScienceID,
		photo,
		roles,
		0,
	)
}

func UserProfileModelToEntity(profile *models.UserProfileModel) *entity2.UserProfileEntity {
	socialModelToEntity := UserSocialsModelToEntity(profile.Socials)

	var photo *string
	if profile.User != nil && profile.User.Photo != "" {
		photo = &profile.User.Photo
	} else if profile.Photo != nil {
		photo = profile.Photo
	}
	var fullName string
	var scienceID string
	var phoneNumber string
	var email string
	var city enum2.CityCode
	var academicDegree enum2.AcademicDegreeCode
	var academicTitle *string
	var orcidID *string
	if profile.User != nil {
		fullName = profile.User.FullName
		scienceID = profile.User.ScienceID
		phoneNumber = profile.User.PhoneNumber
		email = profile.User.Email
		city = profile.User.City
		academicDegree = profile.User.AcademicDegree
		academicTitle = profile.User.AcademicTitle
		orcidID = profile.User.ORCIDID
	}

	return entity2.NewUserProfileEntity(
		profile.ID,
		profile.UserID,
		fullName,
		scienceID,
		email,
		profile.Bio,
		photo,
		phoneNumber,
		socialModelToEntity,
		nil,
		nil,
		&city,
		&academicDegree,
		academicTitle,
		orcidID,
		0,
		0,
	)
}
