package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func DoctorateModelToEntity(model *models.DoctorateModel) *entity.DoctorateEntity {

	return entity.NewDoctorateEntity(
		model.ID,
		model.ExternalID,
		model.DcType,
		model.EduLang,
		model.Status,
		model.StatusCode,
		model.AdmissionYear,
		model.DirectionName,
		model.DirectionCode,
		model.AdvisorFullName,
		model.AdvisorPin,
		model.ScientificWorkName,
		model.OrganizationTin,
		//model.OrganizationID,
		model.UserID,
	)
}

func DoctorateEntityToModel(userID uint, entity *entity.DoctorateEntity) *models.DoctorateModel {
	return models.NewDoctorateModel(
		entity.ExternalID,
		entity.DcType,
		entity.EduLang,
		entity.Status,
		entity.StatusCode,
		entity.AdmissionYear,
		entity.DirectionName,
		entity.DirectionCode,
		entity.AdvisorFullName,
		entity.AdvisorPin,
		entity.ScientificWorkName,
		entity.OrganizationTin,
		//entity.OrganizationID,
		&userID,
	)
}

func DoctorateModelFields() []string {
	return []string{
		"external_id",
		"dc_type",
		"edu_lang",
		"status",
		"status_code",
		"admission_year",
		"direction_name",
		"direction_code",
		"advisor_full_name",
		"advisor_pin",
		"scientific_work_name",
	}
}
