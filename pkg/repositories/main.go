package repositories

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Context struct {
	PlanetRepository     *PlanetRepository
	InviteCodeRepository *InviteCodeRepository
	ConfigRepository     *ConfigRepository
	_db                  *gorm.DB
}

func NewContext(db *gorm.DB) *Context {
	return &Context{
		PlanetRepository:     NewPlanetRepository(db),
		InviteCodeRepository: NewInviteCodeRepository(db),
		ConfigRepository:     NewConfigRepository(db),
		_db:                  db,
	}
}

func (ctx *Context) ExecuteTransaction(transaction func(*gorm.DB) bool) bool {
	err := ctx._db.Transaction(
		func(tx *gorm.DB) error {
			if transaction(ctx._db) {
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
