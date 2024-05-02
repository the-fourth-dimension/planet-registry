package InviteHandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

func (h *inviteHandler) post(ctx *gin.Context) {
	var input struct {
		Code string `json:"code"`
	}

	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.Error(HttpError.NewHttpError("invalid input", err.Error(), http.StatusBadRequest))
		return
	}

	findInviteResult := h.ctx.InviteRepository.FindFirst(&models.Invite{Code: input.Code})
	if findInviteResult.Error != nil {
		if !errors.Is(findInviteResult.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithError(http.StatusInternalServerError, findInviteResult.Error)
			return
		}
	} else {
		ctx.Error(HttpError.NewHttpError("code already exists", input.Code, http.StatusConflict))
		return
	}
	findInviteResult.Result.Code = input.Code
	saveInviteResult := h.ctx.InviteRepository.Save(&findInviteResult.Result)
	if saveInviteResult.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, saveInviteResult.Error)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": saveInviteResult.Result.ID})
}

func (h *inviteHandler) get(ctx *gin.Context) {
	code := ctx.Query("code")
	findObj := models.Invite{}
	if code != "" {
		findObj.Code = code
	}
	findInvitesResult := h.ctx.InviteRepository.Find(&findObj)
	if findInvitesResult.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, findInvitesResult.Error)
		return
	}
	ctx.JSON(http.StatusOK, findInvitesResult.Result)
}
