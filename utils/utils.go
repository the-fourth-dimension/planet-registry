package utils

import (
	"github.com/the_fourth_dimension/planet_registry/models"
	"golang.org/x/crypto/bcrypt"
)

func verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hashedPassword))
}

func LoginCheck(planetId string, password string) error {
	var err error

	planet := models.Planet {}
	err = models.DB.Model(models.Planet {}).Where("planetId = ?", planetId).Take(&planet).Error

	if err != nil {
		return err
	}

	err = verifyPassword(password, planet.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return err
	}

	return nil
}