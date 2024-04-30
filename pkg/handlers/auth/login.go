package handlers_auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/types"
)

func (h *AuthHandler) postLogin(ctx *gin.Context) {
	var input types.Credentials

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(HttpError.NewHttpError("Invalid input", err.Error(), http.StatusBadRequest))
		return
	}

	planet := models.Planet{}

	planet.PlanetId = input.PlanetId

	existingPlanetResult := h.PlanetRepository.FindFirst(&planet)
	if existingPlanetResult.Error != nil {
		if errors.Is(existingPlanetResult.Error, gorm.ErrRecordNotFound) {
			ctx.Error(HttpError.NewHttpError("Not found", "planetId: "+planet.PlanetId, http.StatusNotFound))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, existingPlanetResult.Error)
		return
	}

	err := lib.VerifyPassword(input.Password, existingPlanetResult.Result.Password)

	if err != nil {
		ctx.Error(HttpError.NewHttpError("Invalid credentials", "password", http.StatusForbidden))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"planet": planet,
	})
}
