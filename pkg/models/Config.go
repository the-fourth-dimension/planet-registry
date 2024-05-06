package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Config struct {
	gorm.Model
	InviteOnly bool `gorm:"not null" json:"inviteOnly"`
}

func (c *Config) ToString() string {
	return fmt.Sprintf("------------------ inviteOnly: %t\n", c.InviteOnly)
}
