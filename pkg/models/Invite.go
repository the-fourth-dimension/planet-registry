package models

import "github.com/jinzhu/gorm"

type Invite struct {
	gorm.Model
	Code string `gorm:"size:255;not null;unique" json:"code"`
}
