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

func (r *ConfigRepository) Save(config *models.Config) RepositoryResult[models.Config] {
	err := r.db.Save(config).Error

	return RepositoryResult[models.Config]{Result: config, Error: err}
}

func (r *ConfigRepository) FindFirst(query *models.Config) RepositoryResult[models.Config] {
	var config models.Config
	err := r.db.Find(config, query).Error
	return RepositoryResult[models.Config]{Result: &config, Error: err}
}
