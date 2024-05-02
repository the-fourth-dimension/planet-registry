package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

type PlanetRepository struct {
	db *gorm.DB
}

func NewPlanetRepository(db *gorm.DB) *PlanetRepository {
	return &PlanetRepository{db}
}

func (r *PlanetRepository) Save(data *models.Planet) RepositoryResult[models.Planet] {
	err := r.db.Save(data).Error

	return RepositoryResult[models.Planet]{Result: *data, Error: err}
}

func (r *PlanetRepository) FindAll(query *models.Planet) RepositoryResult[[]models.Planet] {
	var data []models.Planet
	err := r.db.Find(&data, query).Error
	return RepositoryResult[[]models.Planet]{Result: data, Error: err}
}

func (r *PlanetRepository) FindFirst(query *models.Planet) RepositoryResult[models.Planet] {
	var data models.Planet
	err := r.db.Find(&data, query).Error
	return RepositoryResult[models.Planet]{Result: data, Error: err}
}
