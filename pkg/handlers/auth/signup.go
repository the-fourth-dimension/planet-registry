package handlers_auth

import (
	"errors"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

func (h *AuthHandler) postSignUp(ctx *gin.Context) {
	findConfigResult := h.ConfigRepository.FindFirst(&models.Config{})

	if findConfigResult.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, findConfigResult.Error)
		if errors.Is(findConfigResult.Error, gorm.ErrRecordNotFound) {
			log.Fatalln("server config not set!")
		}
		return
	}

	planet := models.Planet{}

	if findConfigResult.Result.InviteOnly {
		var input CredentialsWithCode
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		inviteCode := h.InviteCodeRepository.FindFirst(&models.InviteCode{Code: input.Code})
		if inviteCode.Error != nil {
			if errors.Is(inviteCode.Error, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Incorrect Code",
				})
				return
			}
			ctx.AbortWithError(http.StatusInternalServerError, inviteCode.Error)
			return
		}
		deleteInviteResult := h.InviteCodeRepository.DeleteOneById(inviteCode.Result.ID)
		if deleteInviteResult.Error != nil {
			ctx.AbortWithError(http.StatusInternalServerError, deleteInviteResult.Error)
			return
		}
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
	hashedPassword, err := lib.HashPassword(planet.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	planet.Password = hashedPassword
	planet.PlanetId = html.EscapeString(strings.TrimSpace(planet.PlanetId))
	savePlanetResult := h.PlanetRepository.Save(&planet)

	if savePlanetResult.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, savePlanetResult.Error)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "registration success",
	})
}
