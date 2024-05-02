package AuthHandler

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
)

func (h *authHandler) postSignUp(ctx *gin.Context) {
	findConfigResult := h.ctx.ConfigRepository.FindFirst(&models.Config{})
	if findConfigResult.Error != nil {
		if errors.Is(findConfigResult.Error, gorm.ErrRecordNotFound) {
			log.Fatalln("server config not set!")
		}
		ctx.AbortWithError(http.StatusInternalServerError, findConfigResult.Error)
		return
	}

	planet := models.Planet{}

	if findConfigResult.Result.InviteOnly {
		var input credentialsWithCode
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.Error(HttpError.NewHttpError("invalid input", err.Error(), http.StatusBadRequest))
			return
		}

		success := h.ctx.ExecuteTransaction(func(tx *gorm.DB) bool {
			inviteCode := h.ctx.InviteCodeRepository.FindFirst(&models.InviteCode{Code: input.Code})
			if inviteCode.Error != nil {
				if errors.Is(inviteCode.Error, gorm.ErrRecordNotFound) {
					ctx.Error(HttpError.NewHttpError("invalid code", input.Code, http.StatusBadRequest))
					return false
				}
				ctx.AbortWithError(http.StatusInternalServerError, inviteCode.Error)
				return false
			}
			deleteInviteResult := h.ctx.InviteCodeRepository.DeleteOneById(inviteCode.Result.ID)
			if deleteInviteResult.Error != nil {
				ctx.AbortWithError(http.StatusInternalServerError, deleteInviteResult.Error)
				return false
			}
			if deleteInviteResult.Result <= 0 {
				ctx.Error(HttpError.NewHttpError("invalid code", input.Code, http.StatusBadRequest))
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
		var input credentials
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
	savePlanetResult := h.ctx.PlanetRepository.Save(&planet)

	if savePlanetResult.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, savePlanetResult.Error)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "registration success",
	})
}
