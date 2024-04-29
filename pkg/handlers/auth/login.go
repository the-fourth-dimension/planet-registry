package handlers_auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/types"
)

func (h *AuthHandler) postLogin(ctx *gin.Context) {
	var input types.Credentials

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	planet := models.Planet{}

	planet.PlanetId = input.PlanetId

	existingPlanetResult := h.PlanetRepository.FindFirst(&planet)
	if existingPlanetResult.Error != nil {
		if errors.Is(existingPlanetResult.Error, gorm.ErrRecordNotFound) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, existingPlanetResult.Error)
		return
	}

	err := lib.VerifyPassword(input.Password, existingPlanetResult.Result.Password)

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "incorrect password",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"planet": planet,
	})
}
