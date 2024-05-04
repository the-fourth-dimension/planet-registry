package PlanetHandler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

func (h *planetHandler) post(ctx *gin.Context) {
	findConfigResult := h.ctx.ConfigRepository.FindFirst(&models.Config{})
	if findConfigResult.Error != nil {
		if errors.Is(findConfigResult.Error, gorm.ErrRecordNotFound) {
			log.Fatalln("server config not set!")
		}
		ctx.AbortWithError(http.StatusInternalServerError, findConfigResult.Error)
		return
	}

	planet := models.Planet{}
	var code = ""
	var inviteCodeID uint
	if findConfigResult.Result.InviteOnly {
		var input credentialsWithCode
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.Error(HttpError.NewHttpError("invalid input", err.Error(), http.StatusBadRequest))
			return
		}
		code = input.Code
		inviteCode := h.ctx.InviteRepository.FindFirst(&models.Invite{Code: code})
		if inviteCode.Error != nil {
			if errors.Is(inviteCode.Error, gorm.ErrRecordNotFound) {
				ctx.Error(HttpError.NewHttpError("invalid code", code, http.StatusBadRequest))
				return
			}
			ctx.AbortWithError(http.StatusInternalServerError, inviteCode.Error)
			return
		}
		inviteCodeID = inviteCode.Result.ID
		planet.PlanetId = input.PlanetId
		planet.Password = input.Password
	} else {
		var input credentials
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.Error(HttpError.NewHttpError("invalid input", err.Error(), http.StatusBadRequest))
			return
		}
		planet.PlanetId = input.PlanetId
		planet.Password = input.Password
	}
	findPlanetResult := h.ctx.PlanetRepository.FindFirst(&models.Planet{PlanetId: planet.PlanetId})
	if findPlanetResult.Error != nil {
		if !errors.Is(findPlanetResult.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithError(http.StatusInternalServerError, findPlanetResult.Error)
			return
		}
	} else {
		ctx.Error(HttpError.NewHttpError("planet already exists", planet.PlanetId, http.StatusConflict))
		return
	}

	success := false
	success = h.ctx.ExecuteTransaction(func(d *gorm.DB, txCtx *repositories.Context) bool {
		if findConfigResult.Result.InviteOnly {
			deleteInviteResult := txCtx.InviteRepository.DeleteOneById(inviteCodeID)
			if deleteInviteResult.Error != nil {
				ctx.AbortWithError(http.StatusInternalServerError, deleteInviteResult.Error)
				return false
			}
			if deleteInviteResult.Result <= 0 {
				ctx.Error(HttpError.NewHttpError("invalid code", code, http.StatusBadRequest))
				return false
			}
		}
		savePlanetResult := txCtx.PlanetRepository.Save(&planet)

		if savePlanetResult.Error != nil {
			ctx.AbortWithError(http.StatusInternalServerError, savePlanetResult.Error)
			return false
		}
		return true
	})
	if !success {
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "registration success",
	})
}
