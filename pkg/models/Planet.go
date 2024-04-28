package models

import (
	"github.com/jinzhu/gorm"
)

type Planet struct {
	gorm.Model
	PlanetId string `gorm:"size:255;not null;unique" json:"planetId"`
	Password string `gorm:"size:255;not null" json:"password"`
}
