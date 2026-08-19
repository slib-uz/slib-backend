package entity

type DegreeEntity struct {
	ID               uint
	DegreeTypeID     uint
	DegreeType       string
	Field            *string
	DegreeStatusID   *uint
	DegreeStatusName *string
	ConfirmedDate    *string
	Protocol         string
	UserID           *uint
}

func NewDegreeEntity(id, degreeTypeID uint, degreeType string, field *string, degreeStatusID *uint, degreeStatusName *string, confirmedDate *string, protocol string, userID *uint) *DegreeEntity {
	return &DegreeEntity{
		ID:               id,
		DegreeTypeID:     degreeTypeID,
		DegreeType:       degreeType,
		Field:            field,
		DegreeStatusID:   degreeStatusID,
		DegreeStatusName: degreeStatusName,
		ConfirmedDate:    confirmedDate,
		Protocol:         protocol,
		UserID:           userID,
	}
}
