package models

import (
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
)

type Admin struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
}

func (planet *Admin) BeforeSave(scope *gorm.Scope) (err error) {
	hashedPassword, hashErr := lib.HashPassword(planet.Password)

	if hashErr != nil {
		err = hashErr
		return err
	}
	scope.SetColumn("password", hashedPassword)
	return
}
