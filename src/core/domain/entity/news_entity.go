package entity

import "time"

type NewsEntity struct {
	ID         uint              `json:"id"`
	Title      map[string]string `json:"title"`
	Body       map[string]string `json:"body"`
	CategoryID uint              `json:"category_id"`
	Image      string            `json:"image"`
	ViewsCount uint              `json:"views_count"`
	CreatedAt  time.Time         `json:"created_at"`
}

func NewNewsEntity(ID uint, title map[string]string, body map[string]string, categoryID uint, image string, viewsCount uint, createdAt time.Time) *NewsEntity {
	return &NewsEntity{ID: ID, Title: title, Body: body, CategoryID: categoryID, Image: image, ViewsCount: viewsCount, CreatedAt: createdAt}
}
