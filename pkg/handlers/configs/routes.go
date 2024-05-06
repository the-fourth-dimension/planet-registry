package ConfigsHandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

func (h *configsHandler) get(ctx *gin.Context) {
	findConfigResult := h.ctx.ConfigRepository.FindFirst(&models.Config{})
	if findConfigResult.Error != nil {
		if errors.Is(findConfigResult.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			panic("Server Config not found!")
		}
		ctx.AbortWithError(http.StatusInternalServerError, findConfigResult.Error)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"config": findConfigResult.Result})
}

type putRequestPayload struct {
	InviteOnly *bool `json:"inviteOnly,omitempty"`
}

func (h *configsHandler) put(ctx *gin.Context) {
	findConfigResult := h.ctx.ConfigRepository.FindFirst(&models.Config{})
	if findConfigResult.Error != nil {
		if errors.Is(findConfigResult.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			panic("Server Config not found!")
		}
		ctx.AbortWithError(http.StatusInternalServerError, findConfigResult.Error)
		return
	}
	var input putRequestPayload
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(HttpError.NewHttpError("invalid input", err.Error(), http.StatusBadRequest))
		return
	}
	changed := false
	if input.InviteOnly != nil {
		findConfigResult.Result.InviteOnly = *input.InviteOnly
		changed = true
	}
	if !changed {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "Nothing to change"})
		return
	}
	saveResult := h.ctx.ConfigRepository.Save(&findConfigResult.Result)
	if saveResult.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, saveResult.Error)
		return
	}
	ctx.Status(http.StatusAccepted)
}
