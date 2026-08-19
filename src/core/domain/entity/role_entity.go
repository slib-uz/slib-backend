package entity

import (
	"slib.uz/src/core/domain/entity/enum"
)

type UserRoleEntity struct {
	ID              uint              `json:"id"`
	UserID          uint              `json:"user_id"`
	User            *UserEntity       `json:"user"`
	Role            enum.UserRole     `json:"role"`
	PublisherID     *uint             `json:"publisher_id"`
	JournalID       *uint             `json:"journal_id"`
	InstitutionID   *uint             `json:"institution_id"`
	PublisherName   *string           `json:"publisher_name"`
	JournalName     map[string]string `json:"journal_name"`
	InstitutionName *string           `json:"institution_name"`
	URL             *string           `json:"url"`
}

func NewUserRoleEntity(ID uint, userID uint, user *UserEntity, role enum.UserRole, publisherID *uint, journalID *uint, institutionID *uint, publisherName *string, journalName map[string]string, institutionName *string, URL *string) *UserRoleEntity {
	return &UserRoleEntity{ID: ID, UserID: userID, User: user, Role: role, PublisherID: publisherID, JournalID: journalID, InstitutionID: institutionID, PublisherName: publisherName, JournalName: journalName, InstitutionName: institutionName, URL: URL}
}

func NewUserRoleInputEntity(userID uint, publisherID, journalID, institutionID *uint, role enum.UserRole, url *string) *UserRoleEntity {
	return &UserRoleEntity{
		UserID:        userID,
		PublisherID:   publisherID,
		JournalID:     journalID,
		InstitutionID: institutionID,
		Role:          role,
		URL:           url,
	}
}

type UserRoleWithBasicUserEntity struct {
	ID              uint              `json:"id"`
	UserID          uint              `json:"user_id"`
	User            *UserBasicEntity  `json:"user"`
	Role            enum.UserRole     `json:"role"`
	PublisherID     *uint             `json:"publisher_id"`
	JournalID       *uint             `json:"journal_id"`
	InstitutionID   *uint             `json:"institution_id"`
	PublisherName   *string           `json:"publisher_name"`
	JournalName     map[string]string `json:"journal_name"`
	InstitutionName *string           `json:"institution_name"`
	URL             *string           `json:"url"`
}

func NewUserRoleWithBasicUserEntity(ID uint, userID uint, user *UserBasicEntity, role enum.UserRole, publisherID *uint, journalID *uint, institutionID *uint, publisherName *string, journalName map[string]string, institutionName *string, URL *string) *UserRoleWithBasicUserEntity {
	return &UserRoleWithBasicUserEntity{ID: ID, UserID: userID, User: user, Role: role, PublisherID: publisherID, JournalID: journalID, InstitutionID: institutionID, PublisherName: publisherName, JournalName: journalName, InstitutionName: institutionName, URL: URL}
}
