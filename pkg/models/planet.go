package models

import (
	"html"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Planet struct {
	gorm.Model
	PlanetId string `gorm:"size:255;not null;unique" json:"planetId"`
	Password string `gorm:"size:255;not null" json:"password"`
}

func (planet *Planet) Save() (*Planet, error) {
	err := DB.Create(&planet).Error
	if err != nil {
		return &Planet{}, err
	}
	return planet, nil
}

func (planet *Planet) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(planet.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	planet.Password = string(hashedPassword)
	planet.PlanetId = html.EscapeString(strings.TrimSpace(planet.PlanetId))

	return nil
}
