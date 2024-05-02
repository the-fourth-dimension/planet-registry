package InviteHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type inviteHandler struct {
	router *gin.Engine
	ctx    *repositories.Context
}

func (h *inviteHandler) RegisterRouter() {
	invite := h.router.Group("/invite")
}

func New(router *gin.Engine, ctx *repositories.Context) *inviteHandler {
	return &inviteHandler{
		router: router,
		ctx:    ctx,
	}
}
