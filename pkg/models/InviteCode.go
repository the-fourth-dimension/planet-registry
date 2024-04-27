package models

import "github.com/jinzhu/gorm"

type InviteCode struct {
	gorm.Model
	Code string `gorm:"size:255;not null;unique" json:"code"`
}

func (inviteCode *InviteCode) Save() (*InviteCode, error) {
	err := DB.Create(&inviteCode).Error
	if err != nil {
		return &InviteCode{}, err
	}
	return inviteCode, nil
}

func (invite *InviteCode) GetAll() ([]*InviteCode, error) {
	invites := []*InviteCode{}
	err := DB.Find(invites).Error
	if err != nil {
		return []*InviteCode{}, err
	}
	return invites, nil
}

func (invite *InviteCode) DeleteOne() error {
	err := DB.Delete(invite).Error
	return err
}
