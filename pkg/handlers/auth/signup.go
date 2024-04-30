package handlers_auth

import (
	"errors"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/types"
)

func (h *AuthHandler) postSignUp(ctx *gin.Context) {
	findConfigResult := h.ConfigRepository.FindFirst(&models.Config{})

	if findConfigResult.Error != nil {
		if errors.Is(findConfigResult.Error, gorm.ErrRecordNotFound) {
			log.Fatalln("server config not set!")
		}
		ctx.AbortWithError(http.StatusInternalServerError, findConfigResult.Error)
		return
	}

	planet := models.Planet{}

	if findConfigResult.Result.InviteOnly {
		var input types.CredentialsWithCode
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.Error(HttpError.NewHttpError("invalid input", err.Error(), http.StatusBadRequest))
			return
		}

		success := h.ExecuteTransaction(func(tx *gorm.DB) bool {
			inviteCode := h.InviteCodeRepository.FindFirst(&models.InviteCode{Code: input.Code})
			if inviteCode.Error != nil {
				if errors.Is(inviteCode.Error, gorm.ErrRecordNotFound) {
					ctx.Error(HttpError.NewHttpError("Invalid Code", input.Code, http.StatusBadRequest))
					return false
				}
				ctx.AbortWithError(http.StatusInternalServerError, inviteCode.Error)
				return false
			}
			deleteInviteResult := h.InviteCodeRepository.DeleteOneById(inviteCode.Result.ID)
			if deleteInviteResult.Error != nil {
				ctx.AbortWithError(http.StatusInternalServerError, deleteInviteResult.Error)
				return false
			}
			if err := ctx.ShouldBindJSON(&input); err != nil {
				ctx.Error(HttpError.NewHttpError("invalid input", err.Error(), http.StatusBadRequest))
				return false
			}
			planet.PlanetId = input.PlanetId
			planet.Password = input.Password
			return true
		})

		if !success {
			return
		}
	} else {
		var input types.Credentials
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.Error(HttpError.NewHttpError("invalid input", err.Error(), http.StatusBadRequest))
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
