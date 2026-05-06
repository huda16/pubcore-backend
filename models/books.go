package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

// Books represents a book entity with Author and Publisher relationships.
type Books struct {
	GormModel
	Title       string      `gorm:"not null" json:"title" valid:"required~Title is required"`
	ISBN        string      `gorm:"uniqueIndex" json:"isbn"`
	Description string      `json:"description"`
	Genre       string      `json:"genre"`
	Year        int         `json:"year"`
	AuthorID    uint        `gorm:"not null" json:"author_id" valid:"required~Author is required"`
	Author      *Authors    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	PublisherID uint        `gorm:"not null" json:"publisher_id" valid:"required~Publisher is required"`
	Publisher   *Publishers `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
}

func (b *Books) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(b)
	if errCreate != nil {
		return errCreate
	}
	return nil
}

func (b *Books) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errUpdate := govalidator.ValidateStruct(b)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}
