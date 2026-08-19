package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type NewsCreateRequest struct {
	Title      map[string]string `json:"title"`
	Body       map[string]string `json:"body"`
	Image      string            `json:"image"`
	CategoryID uint              `json:"category_id"`
}

func NewsCreateRequestToEntity(request *NewsCreateRequest) *entity.NewsEntity {
	return entity.NewNewsEntity(0, request.Title, request.Body, request.CategoryID, request.Image, 0, time.Now())
}
