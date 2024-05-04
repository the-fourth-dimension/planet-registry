package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

type Invite struct {
	gorm.Model
	Code string `gorm:"size:255;not null;unique" json:"code"`
}

func (i *Invite) BeforeSave(scope *gorm.Scope) (err error) {
	code := strings.Trim(i.Code, " ")
	if len(code) < 4 {
		err = errors.New("code-length-too-short")
	}
	scope.SetColumn("code", code)
	return
}
