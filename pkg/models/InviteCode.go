package models

import "github.com/jinzhu/gorm"

type InviteCode struct {
	gorm.Model
	Code string `gorm:"size:255;not null;unique" json:"code"`
}
