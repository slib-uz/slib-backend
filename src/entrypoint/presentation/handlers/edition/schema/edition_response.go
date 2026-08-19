package schema

import (
	"slib.uz/src/core/domain/entity"
)

type EditionResponse struct {
	Year     int                    `json:"year"`
	Editions []entity.EditionEntity `json:"editions"`
}

type EditionListResponse struct {
	Years   []EditionResponse `json:"years"`
	Total   int               `json:"total"`
	Page    int               `json:"page"`
	PerPage int               `json:"per_page"`
}
