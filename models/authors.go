package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

// Authors represents an author entity in the publishing platform.
type Authors struct {
	GormModel
	Name        string  `gorm:"not null" json:"name" valid:"required~Author name is required"`
	Bio         string  `json:"bio"`
	Nationality string  `json:"nationality"`
	Books       []Books `gorm:"foreignKey:AuthorID" json:"books,omitempty"`
}

func (a *Authors) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(a)
	if errCreate != nil {
		return errCreate
	}
	return nil
}

func (a *Authors) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errUpdate := govalidator.ValidateStruct(a)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}
