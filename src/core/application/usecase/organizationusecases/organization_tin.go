package organizationusecases

import (
	"slib.uz/src/core/application/response"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/utils"
)

const organizationAlreadyExistsMessage = "Bu tashkilot allaqachon kiritilgan. Iltimos, uni 'Tashkilotlar' bo'limidan qidirib toping."

func normalizeAndValidateOrganizationTin(tin *string) (string, error) {
	raw := ""
	if tin != nil {
		raw = *tin
	}
	normalized := utils.NormalizeTin(raw)
	if !utils.IsValidOrganizationTin(normalized) {
		return "", response.NewFailResponse(400, "STIR 9 xonali raqam bo'lishi kerak")
	}
	return normalized, nil
}

func organizationAlreadyExists(existing *entity2.OrganizationEntity) error {
	return response.NewOptionalResponse(409, response.CodeOrganizationAlreadyExists, existing, organizationAlreadyExistsMessage)
}
