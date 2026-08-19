package mapper

import (
	"slib.uz/src/core/domain/entity"
)

func PublisherEntityToROIPublisherEntity(publisher *entity.PublisherEntity) *entity.ROIPublisherEntity {
	var name, shortName string

	if publisher.Name != nil {
		name = *publisher.Name
	}
	if publisher.ShortName != nil {
		shortName = *publisher.ShortName
	}

	return entity.NewROIPublisherEntity(
		name,
		shortName,
		publisher.Tin,
	)
}
