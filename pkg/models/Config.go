package models

import "github.com/jinzhu/gorm"

type Config struct {
	gorm.Model
	InviteOnly bool `gorm:"not null" json:"inviteOnly"`
}
