package InviteHandler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/errors/HttpError"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

func (h *inviteHandler) post(ctx *gin.Context) {
	var input struct {
		Code string `json:"code"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
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

func (h *inviteHandler) deleteById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(HttpError.NewHttpError("Missing param fields", "id", http.StatusBadRequest))
		return
	}
	uintId, error := strconv.ParseUint(id, 10, 64)
	if error != nil {
		ctx.Error(HttpError.NewHttpError("Invalid field", "id", http.StatusBadRequest))
	}
	findResult := h.ctx.InviteRepository.FindFirst(&models.Invite{Model: gorm.Model{ID: uint(uintId)}})
	if findResult.Error != nil {
		if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
			ctx.Error(HttpError.NewHttpError("Invite.ID Not found", id, http.StatusNotFound))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, findResult.Error)
		return
	}
	deleteResult := h.ctx.InviteRepository.DeleteOneById(findResult.Result.ID)
	if deleteResult.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, deleteResult.Error)
		return
	}
	ctx.Status(http.StatusNoContent)
}

type putRequestPayload struct {
	Code string `json:"code"`
}

func (h *inviteHandler) putById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(HttpError.NewHttpError("missing param fields", "id", http.StatusBadRequest))
		return
	}
	uintId, error := strconv.ParseUint(id, 10, 64)
	if error != nil {
		ctx.Error(HttpError.NewHttpError("invalid field", "id", http.StatusBadRequest))
	}
	findResult := h.ctx.InviteRepository.FindFirst(&models.Invite{Model: gorm.Model{ID: uint(uintId)}})
	if findResult.Error != nil {
		if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
			ctx.Error(HttpError.NewHttpError("Invite.ID Not found", id, http.StatusNotFound))
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
	if input.Code != findResult.Result.Code {
		findResult.Result.Code = input.Code
		changed = true
	}
	if changed {
		saveResult := h.ctx.InviteRepository.Save(&findResult.Result)
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
