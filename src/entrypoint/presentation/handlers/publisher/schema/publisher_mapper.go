package schema

import (
	"slib.uz/src/core/domain/entity"
)

func PublisherResToDto(res *PublisherCreateRequest) *entity.PublisherEntity {
	publisher := entity.NewPublisherEntity(
		res.ID,
		res.Tin,
		res.Name,
		res.ShortName,
		res.Logo,
		res.PhoneNumber,
		res.FaxNumber,
		res.Email,
		res.Website,
		res.Address,
		res.Description,
		res.IsActive,
	)
	publisher.RegionID = res.RegionID
	publisher.DistrictID = res.DistrictID
	return publisher
}
