package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/entrypoint/presentation/handlers/job/schema"
)

func JobCreateUpdateRequestToDTO(request *schema.JobCreateUpdateRequest) *entity2.JobCreateUpdateEntity {
	return entity2.NewJobCreateUpdateEntity(request.PositionName, request.OrganizationID)
}

func JobDTOToCreateUpdateResponse(jobDTO *entity2.JobEntity) *schema.JobCreateUpdateResponse {
	return &schema.JobCreateUpdateResponse{
		ID:               jobDTO.ID,
		OrganizationTin:  jobDTO.OrganizationTin,
		OrganizationID:   jobDTO.OrganizationID,
		OrganizationName: jobDTO.OrganizationName,
		UserID:           jobDTO.UserID,
		PositionName:     jobDTO.PositionName,
	}
}
