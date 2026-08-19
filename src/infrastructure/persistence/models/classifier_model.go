package models

import "gorm.io/gorm"

type SoatoClassifierModel struct {
	gorm.Model

	Code   uint    `gorm:"uniqueIndex;not null'"`
	NameUz *string `gorm:"size:255"`
	NameRu *string `gorm:"size:255"`
	NameEn *string `gorm:"size:255"`

	ParentID *uint                  `gorm:"index"`
	Parent   *SoatoClassifierModel  `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL"` // many2one
	Children []SoatoClassifierModel `gorm:"foreignKey:ParentID"`                              // one2many
}

func (this SoatoClassifierModel) Name() map[string]string {
	return map[string]string{
		"uz": getStringValue(this.NameUz),
		"ru": getStringValue(this.NameRu),
		"en": getStringValue(this.NameEn),
	}
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (SoatoClassifierModel) TableName() string {
	return "soato_classifiers"
}

func NewSoatoClassifierModel(code uint, parentID *uint, nameUz, nameRu, nameEn *string) *SoatoClassifierModel {
	return &SoatoClassifierModel{
		Code:     code,
		ParentID: parentID,
		NameUz:   nameUz,
		NameRu:   nameRu,
		NameEn:   nameEn,
	}
}
