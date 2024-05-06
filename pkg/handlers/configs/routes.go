package ConfigsHandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
