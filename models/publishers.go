package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

// Publishers represents a publisher entity in the publishing platform.
type Publishers struct {
	GormModel
	Name    string  `gorm:"not null;uniqueIndex" json:"name" valid:"required~Publisher name is required"`
	Address string  `json:"address"`
	Website string  `json:"website"`
	Books   []Books `gorm:"foreignKey:PublisherID" json:"books,omitempty"`
}

func (p *Publishers) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)
	if errCreate != nil {
		return errCreate
	}
	return nil
}

func (p *Publishers) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errUpdate := govalidator.ValidateStruct(p)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}
