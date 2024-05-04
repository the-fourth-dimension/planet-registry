package models

import (
	"html"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
)

type Planet struct {
	gorm.Model
	PlanetId string `gorm:"size:255;not null;unique" json:"planetId"`
	Password string `gorm:"size:255;not null" json:"password"`
}

func (planet *Planet) BeforeSave(scope *gorm.Scope) (err error) {
	hashedPassword, hashErr := lib.HashPassword(planet.Password)

	if hashErr != nil {
		err = hashErr
		return err
	}
	scope.SetColumn("password", hashedPassword)
	planet.PlanetId = html.EscapeString(strings.TrimSpace(planet.PlanetId))
	return
}
