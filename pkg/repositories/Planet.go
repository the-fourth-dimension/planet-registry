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

func (r *PlanetRepository) Save(planet *models.Planet) RepositoryResult[models.Planet] {
	err := r.db.Save(planet).Error

	return RepositoryResult[models.Planet]{Result: planet, Error: err}
}

func (r *PlanetRepository) FindAll(query *models.Planet) RepositoryResult[[]models.Planet] {
	var planets []models.Planet
	err := r.db.Find(&planets, query).Error
	return RepositoryResult[[]models.Planet]{Result: &planets, Error: err}
}

func (r *PlanetRepository) FindFirst(query *models.Planet) RepositoryResult[models.Planet] {
	var planet models.Planet
	err := r.db.Find(planet, query).Error
	return RepositoryResult[models.Planet]{Result: &planet, Error: err}
}
