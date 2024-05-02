package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

type InviteRepository struct {
	db *gorm.DB
}

func NewInviteCodeRepository(db *gorm.DB) *InviteRepository {
	return &InviteRepository{db}
}

func (r *InviteRepository) Save(data *models.Invite) RepositoryResult[models.Invite] {
	err := r.db.Save(data).Error

	return RepositoryResult[models.Invite]{Result: *data, Error: err}
}

func (r *InviteRepository) Find(query *models.Invite) RepositoryResult[[]models.Invite] {
	var data []models.Invite
	err := r.db.Find(&data, query).Error
	return RepositoryResult[[]models.Invite]{Result: data, Error: err}
}

func (r *InviteRepository) FindFirst(query *models.Invite) RepositoryResult[models.Invite] {
	var data models.Invite
	err := r.db.Find(&data, query).Error
	return RepositoryResult[models.Invite]{Result: data, Error: err}
}

func (r *InviteRepository) DeleteOneById(ID uint) RepositoryResult[int64] {
	result := r.db.Delete(&models.Invite{}, ID)
	return RepositoryResult[int64]{Error: result.Error, Result: result.RowsAffected}
}
