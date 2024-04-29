package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

type InviteCodeRepository struct {
	db *gorm.DB
}

func NewInviteCodeRepository(db *gorm.DB) *InviteCodeRepository {
	return &InviteCodeRepository{db}
}

func (r *InviteCodeRepository) Save(data *models.InviteCode) RepositoryResult[models.InviteCode] {
	err := r.db.Save(data).Error

	return RepositoryResult[models.InviteCode]{Result: data, Error: err}
}

func (r *InviteCodeRepository) FindFirst(query *models.InviteCode) RepositoryResult[models.InviteCode] {
	var data models.InviteCode
	err := r.db.Find(&data, query).Error
	return RepositoryResult[models.InviteCode]{Result: &data, Error: err}
}

func (r *InviteCodeRepository) DeleteOneById(ID uint) RepositoryResult[any] {
	err := r.db.Delete(&models.InviteCode{}, ID).Error
	return RepositoryResult[any]{Error: err}
}
