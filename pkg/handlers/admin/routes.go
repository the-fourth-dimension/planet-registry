package AdminHandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

func (h *adminHandler) get(ctx *gin.Context) {
	findResult := h.ctx.AdminRepository.Find(&models.Admin{})
	if findResult.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, findResult.Error)
	}
	ctx.JSON(http.StatusOK, findResult)
}

func (h *adminHandler) post(ctx *gin.Context) {
	var input credentials

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(HttpError.NewHttpError("Invalid input", err.Error(), http.StatusBadRequest))
		return
	}

	adminFindByUsernameResult := h.ctx.AdminRepository.FindFirst(&models.Admin{Username: input.username})
	if adminFindByUsernameResult.Error == nil {
		ctx.Error(HttpError.NewHttpError("Username already exists", input.username, http.StatusConflict))
		return
	}
	if adminFindByUsernameResult.Error != nil {
		if !errors.Is(adminFindByUsernameResult.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithError(500, adminFindByUsernameResult.Error)
			return
		}
	}

	admin := models.Admin{Username: input.username, Password: input.password}
	saveAdminResult := h.ctx.AdminRepository.Save(&admin)
	if saveAdminResult.Error != nil {
		ctx.AbortWithError(500, saveAdminResult.Error)
	}

	ctx.Status(http.StatusCreated)
	return
}
