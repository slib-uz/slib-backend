package mapper

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/gateway/response"
)

func UserDataResToEntity(data *response.UserDataResponse) *entity.UserEntity {
	var birthDate *time.Time
	parsed, err := time.Parse("2006-01-02", data.BirthDate)
	if err == nil && !parsed.IsZero() {
		birthDate = &parsed
	}
	return entity.NewUserEntity(
		0,
		data.ScienceID,
		&data.Pin,
		data.FirstName,
		data.LastName,
		data.MiddleName,
		data.FullName,
		data.PhoneNumber,
		false,
		data.Gender,
		birthDate,
		nil,
		"",
		data.Email,
		"",
		"",
		nil,
		nil,
		0, // ArticleCount
	)
}
