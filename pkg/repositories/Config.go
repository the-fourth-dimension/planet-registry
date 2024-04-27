package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

type ConfigRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) *ConfigRepository {
	return &ConfigRepository{db}
}

func (r *ConfigRepository) Save(data *models.Config) RepositoryResult[models.Config] {
	err := r.db.Save(data).Error

	return RepositoryResult[models.Config]{Result: data, Error: err}
}

func (r *ConfigRepository) FindFirst(query *models.Config) RepositoryResult[models.Config] {
	var data models.Config
	err := r.db.Find(data, query).Error
	return RepositoryResult[models.Config]{Result: &data, Error: err}
}
