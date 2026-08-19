package models

type TagModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:32;not null;uniqueIndex:idx_tags_lang_name"`
	Lang string `gorm:"size:5;not null;default:uz;uniqueIndex:idx_tags_lang_name"`
}

func NewTagModel(id uint, name, lang string) *TagModel {
	return &TagModel{ID: id, Name: name, Lang: lang}
}

func (*TagModel) TableName() string {
	return "tags"
}
