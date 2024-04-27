package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/utils"
)

type Credentials struct {
	PlanetId string `json:"planetId" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CredentialsWithCode struct {
	Credentials
	Code string `json:"code" binding:"required"`
}

func Signup(ctx *gin.Context) {
	var input CredentialsWithCode

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	planet := models.Planet{}

	planet.PlanetId = input.PlanetId
	planet.Password = input.Password

	planet.BeforeSave()
	_, err := planet.Save()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "registration success",
	})

}

func Login(ctx *gin.Context) {
	var input Credentials

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	planet := models.Planet{}

	planet.PlanetId = input.PlanetId
	planet.Password = input.Password

	err := utils.LoginCheck(planet.PlanetId, planet.Password)

	if err != nil {
		ctx.Status(403)
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"planet": planet,
	})
}
