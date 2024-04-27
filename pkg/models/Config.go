package models

import "github.com/jinzhu/gorm"

type Config struct {
	gorm.Model
	InviteOnly bool `gorm:"not null" json:"inviteOnly"`
}

func GetConfig() (*Config, error) {
	config := Config{}
	err := DB.First(&config).Error
	return &config, err
}
