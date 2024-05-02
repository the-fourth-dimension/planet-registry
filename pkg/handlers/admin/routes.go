package AdminHandler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

func (h *adminHandler) get(ctx *gin.Context) {
	findResult := h.ctx.AdminRepository.Find(&models.Admin{})
	if findResult.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, findResult.Error)
	}
	ctx.JSON(http.StatusOK, findResult)
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
		hashed, err := lib.HashPassword(input.Password)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if hashed != findResult.Result.Password {
			findResult.Result.Password = hashed
			changed = true
		}
	}
	if input.Username != "" && input.Password != findResult.Result.Password {
		findResult.Result.Username = input.Username
		changed = true
	}
	if changed {
		saveResult := h.ctx.AdminRepository.Save(&findResult.Result)
		if saveResult.Error != nil {
			ctx.AbortWithError(http.StatusInternalServerError, saveResult.Error)
			return
		}
	}
	if !changed {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "Nothing to change"})
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
	hashedPassword, err := lib.HashPassword(input.password)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	admin := models.Admin{Username: input.username, Password: hashedPassword}
	saveAdminResult := h.ctx.AdminRepository.Save(&admin)
	if saveAdminResult.Error != nil {
		ctx.AbortWithError(500, saveAdminResult.Error)
	}

	ctx.Status(http.StatusCreated)
	return
}
