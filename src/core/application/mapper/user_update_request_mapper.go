package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/entrypoint/presentation/handlers/profile/schema"
)

func UserUpdateRequestToEntity(request *schema.UserUpdateRequest) *entity.UserUpdateEntity {
	return entity.NewUserUpdateEntity(
		request.Photo,
		request.Email,
		request.AcademicDegree,
		request.AcademicTitle,
		request.ORCIDID,
	)
}

func UserUpdateEntityToResponse(entity *entity.UserUpdateEntity) *schema.UserUpdateResponse {
	return schema.NewUserUpdateResponse(
		entity.Photo,
		entity.Email,
		entity.AcademicDegree,
		entity.AcademicTitle,
		entity.ORCIDID,
	)
}
