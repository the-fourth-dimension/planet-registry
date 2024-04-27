package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/models"
	"github.com/the_fourth_dimension/planet_registry/utils"
)

type SignupInput struct {
	PlanetId string `json:"planetId" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type LoginInput struct {
	PlanetId string `json:"planetId" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Signup(ctx *gin.Context) {
	var input SignupInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	planet := models.Planet{}

	planet.PlanetId = input.PlanetId
	planet.Password = input.Password

	_, err := planet.SavePlanet()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return 
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "registration success",
	})

}

func Login(ctx *gin.Context) {
	var input LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error(),
		})
		return
	}

	planet := models.Planet{}

	planet.PlanetId = input.PlanetId
	planet.Password = input.Password

	err := utils.LoginCheck(planet.PlanetId, planet.Password)
	
	if err != nil  {
		ctx.Status(403)
		ctx.JSON(http.StatusNotFound, gin.H {
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H {
		"planet": planet, 
	})	
}