package AdminHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type adminHandler struct {
	router *gin.Engine
	ctx    *repositories.Context
}

func (h *adminHandler) RegisterRouter() {
	admin := h.router.Group("/admin")
}

func New(router *gin.Engine, ctx *repositories.Context) *adminHandler {
	return &adminHandler{
		router: router,
		ctx:    ctx,
	}
}
