package AdminHandler

import (
	"errors"
	"net/http"
	"strconv"

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
	ctx.JSON(http.StatusOK, findResult.Result)
}

type putRequestPayload struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (h *adminHandler) putById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(HttpError.NewHttpError("Missing param fields", "id", http.StatusBadRequest))
		return
	}
	uintId, error := strconv.ParseUint(id, 10, 64)
	if error != nil {
		ctx.Error(HttpError.NewHttpError("Invalid field", "id", http.StatusBadRequest))
	}
	findResult := h.ctx.AdminRepository.FindFirst(&models.Admin{Model: gorm.Model{ID: uint(uintId)}})
	if findResult.Error != nil {
		if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
			ctx.Error(HttpError.NewHttpError("Admin.ID Not found", id, http.StatusNotFound))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, findResult.Error)
		return
	}
	var input putRequestPayload
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(HttpError.NewHttpError("invalid input", err.Error(), http.StatusBadRequest))
		return
	}
	changed := false
	if input.Password != "" {
		findResult.Result.Password = input.Password
		changed = true
	}
	if input.Username != "" {
		findResult.Result.Username = input.Username
		changed = true
	}
	if !changed {
		ctx.Status(http.StatusNoContent)
		return
	}
	saveResult := h.ctx.AdminRepository.Save(&findResult.Result)
	if saveResult.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, saveResult.Error)
		return
	}

	ctx.Status(http.StatusAccepted)
}

func (h *adminHandler) post(ctx *gin.Context) {
	var input credentials

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(HttpError.NewHttpError("Invalid input", err.Error(), http.StatusBadRequest))
		return
	}

	adminFindByUsernameResult := h.ctx.AdminRepository.FindFirst(&models.Admin{Username: input.Username})
	if adminFindByUsernameResult.Error == nil {
		ctx.Error(HttpError.NewHttpError("Username already exists", input.Username, http.StatusConflict))
		return
	}
	if adminFindByUsernameResult.Error != nil {
		if !errors.Is(adminFindByUsernameResult.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithError(500, adminFindByUsernameResult.Error)
			return
		}
	}

	admin := models.Admin{Username: input.Username, Password: input.Password}
	saveAdminResult := h.ctx.AdminRepository.Save(&admin)
	if saveAdminResult.Error != nil {
		ctx.AbortWithError(500, saveAdminResult.Error)
	}

	ctx.Status(http.StatusCreated)
	return
}

func (h *adminHandler) deleteById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(HttpError.NewHttpError("missing param fields", "id", http.StatusBadRequest))
		return
	}
	uintId, error := strconv.ParseUint(id, 10, 64)
	if error != nil {
		ctx.Error(HttpError.NewHttpError("invalid field", "id", http.StatusBadRequest))
	}
	deleteResult := h.ctx.AdminRepository.DeleteOneById(uint(uintId))
	if deleteResult.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, deleteResult.Error)
		return
	}
	if deleteResult.Result <= 0 {
		ctx.Error(HttpError.NewHttpError("record not found", id, http.StatusNotFound))
		return
	}
	ctx.Status(http.StatusNoContent)
	return
}
