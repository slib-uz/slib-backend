package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/gateway/response"
)

func AuthorResToEntity(authorRes *response.AuthorResponse) *entity.AuthorEntity {
	return entity.NewAuthorEntity(
		0,
		authorRes.FullName,
		authorRes.ScienceID,
		0,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
}

func AuthorEntityToRes(author *entity.AuthorEntity) *response.AuthorResponse {
	return response.NewAuthorResponse(author.ID, author.FullName, author.ScienceID)
}
