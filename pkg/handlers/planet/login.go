package PlanetHandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

func (h *planetHandler) postLogin(ctx *gin.Context) {
	var input credentials

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(HttpError.NewHttpError("invalid input", err.Error(), http.StatusBadRequest))
		return
	}

	planet := models.Planet{}

	planet.PlanetId = input.PlanetId

	existingPlanetResult := h.ctx.PlanetRepository.FindFirst(&planet)
	if existingPlanetResult.Error != nil {
		if errors.Is(existingPlanetResult.Error, gorm.ErrRecordNotFound) {
			ctx.Error(HttpError.NewHttpError("not found", "planetId: "+planet.PlanetId, http.StatusNotFound))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, existingPlanetResult.Error)
		return
	}

	err := lib.VerifyPassword(input.Password, existingPlanetResult.Result.Password)

	if err != nil {
		ctx.Error(HttpError.NewHttpError("invalid credentials", "password", http.StatusForbidden))
		return
	}

	ctx.Status(200)
}
