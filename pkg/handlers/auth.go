package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

// TODO
// -> USE Transaction in signup
// -> separate methods for models into repositories
func Signup(ctx *gin.Context) {

	config, err := models.GetConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatalln("server config not set!")
		}
		return
	}

	planet := models.Planet{}
	if config.InviteOnly {
		var input CredentialsWithCode
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		inviteCode := models.InviteCode{Code: input.Code}
		_, err := inviteCode.FindOne()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Incorrect Code",
				})
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}
		err = inviteCode.DeleteOne()
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		planet.PlanetId = input.PlanetId
		planet.Password = input.Password
	} else {
		var input Credentials
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		planet.PlanetId = input.PlanetId
		planet.Password = input.Password
	}

	err = planet.BeforeSave()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = planet.Save()

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
