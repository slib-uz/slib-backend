package mapper

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func JournalUpdateEntityToModel(oldJournal *models.JournalModel, e *entity.JournalCreateEntity) models.JournalModel {
	oldJournal.Name = ToGormJson(e.Name)
	oldJournal.ShortName = ToGormJson(e.ShortName)
	oldJournal.Description = ToGormJson(e.Description)
	oldJournal.Rule = ToGormJson(e.Rule)
	oldJournal.ArticlePublishConditions = ToGormJson(e.ArticlePublishConditions)
	oldJournal.Address = ToGormJson(e.Address)
	oldJournal.ISSNPaper = e.ISSNPaper
	oldJournal.ISSNOnline = e.ISSNOnline
	oldJournal.Website = e.Website
	oldJournal.CertificateFile = e.CertificateFile
	oldJournal.Email = e.Email
	oldJournal.PhoneNumber = e.PhoneNumber
	oldJournal.CoverImageFile = e.CoverImageFile
	// Parse date string to time.Time
	if e.DateOfEstablished != "" {
		establishedDate, _ := time.Parse("2006-01-02", e.DateOfEstablished)
		oldJournal.EstablishedDate = establishedDate
	}
	oldJournal.PublishingPrice = e.PublishingPrice
	oldJournal.SellingPrice = e.SellingPrice
	oldJournal.OAKCertificateFile = e.OAKCertificateFile
	oldJournal.PeerReviewMethod = e.PeerReviewMethod
	oldJournal.AccessType = e.AccessType
	oldJournal.SubmissionAccess = e.SubmissionAccess
	oldJournal.CommentAccess = e.CommentAccess
	oldJournal.SocialNetworks = ToGormJson(e.SocialNetworks)
	oldJournal.SupportLink = e.SupportLink
	oldJournal.RegionID = e.RegionID
	oldJournal.DistrictID = e.DistrictID
	return *oldJournal
}
