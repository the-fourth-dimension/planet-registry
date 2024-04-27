package models

import "github.com/jinzhu/gorm"

type InviteCode struct {
	gorm.Model
	Code string `gorm:"size:255;not null;unique" json:"code"`
}

func (inviteCode *InviteCode) Save() (*InviteCode, error) {
	err := DB.Create(&inviteCode).Error
	return inviteCode, err
}

func (invite *InviteCode) GetAll() ([]*InviteCode, error) {
	invites := []*InviteCode{}
	err := DB.Find(invites).Error
	return invites, err
}

func (invite *InviteCode) FindOne() (*InviteCode, error) {
	err := DB.First(invite).Error
	return invite, err
}

func (invite *InviteCode) DeleteOne() error {
	err := DB.Delete(invite).Error
	return err
}
