package models

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type ReportModel struct {
	gorm.Model

	Reason     string          `gorm:"size:512;not null"`
	TargetID   uint            `gorm:"not null"`
	TargetType enum.ReportType `gorm:"not null"`
	ReporterID uint            `gorm:"not null"`
	Reporter   *UserModel      `gorm:"foreignKey:ReporterID"`
	Target     interface{}     `gorm:"-" json:"target,omitempty"`
}

func NewReportModel(reason string, targetID uint, targetType enum.ReportType, reporterID uint) *ReportModel {
	return &ReportModel{Reason: reason, TargetID: targetID, TargetType: targetType, ReporterID: reporterID}
}

func (ReportModel) TableName() string { return "report" }
