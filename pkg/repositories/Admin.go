package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db}
}

func (r *AdminRepository) Save(data *models.Admin) RepositoryResult[models.Admin] {
	err := r.db.Save(data).Error

	return RepositoryResult[models.Admin]{Result: *data, Error: err}
}

func (r *AdminRepository) FindFirst(query *models.Admin) RepositoryResult[models.Admin] {
	var data models.Admin
	err := r.db.First(&data, query).Error
	return RepositoryResult[models.Admin]{Result: data, Error: err}
}

func (r *AdminRepository) Find(query *models.Admin) RepositoryResult[[]models.Admin] {
	var data []models.Admin
	err := r.db.Find(&data, query).Error
	return RepositoryResult[[]models.Admin]{Result: data, Error: err}
}

func (r *AdminRepository) DeleteOneById(ID uint) RepositoryResult[int64] {
	result := r.db.Unscoped().Delete(&models.Admin{}, ID)
	return RepositoryResult[int64]{Result: result.RowsAffected, Error: result.Error}
}
