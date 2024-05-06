package repositories

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
)

type Context struct {
	PlanetRepository *PlanetRepository
	InviteRepository *InviteRepository
	ConfigRepository *ConfigRepository
	AdminRepository  *AdminRepository
	db               *gorm.DB
	logger           *log.Logger
}

func NewContext(db *gorm.DB, logger *log.Logger) *Context {
	return &Context{
		PlanetRepository: NewPlanetRepository(db),
		InviteRepository: NewInviteCodeRepository(db),
		ConfigRepository: NewConfigRepository(db),
		AdminRepository:  NewAdminRepository(db),
		db:               db,
		logger:           logger,
	}
}

func (ctx *Context) ExecuteTransaction(transaction func(*gorm.DB, *Context) bool) bool {
	err := ctx.db.Transaction(
		func(tx *gorm.DB) error {
			if transaction(tx, NewContext(tx, ctx.logger)) {
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
