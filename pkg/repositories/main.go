package repositories

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Context struct {
	PlanetRepository     *PlanetRepository
	InviteCodeRepository *InviteCodeRepository
	ConfigRepository     *ConfigRepository
	AdminRepository      *AdminRepository
	db                   *gorm.DB
}

func NewContext(db *gorm.DB) *Context {
	return &Context{
		PlanetRepository:     NewPlanetRepository(db),
		InviteCodeRepository: NewInviteCodeRepository(db),
		ConfigRepository:     NewConfigRepository(db),
		AdminRepository:      NewAdminRepository(db),
		db:                   db,
	}
}

func (ctx *Context) ExecuteTransaction(transaction func(*gorm.DB) bool) bool {
	err := ctx.db.Transaction(
		func(tx *gorm.DB) error {
			if transaction(ctx.db) {
				return nil
			}
			return errors.New("Transaction Failed")
		},
	)
	if err != nil {
		return false
	}
	return true
}
