package InvitesHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/middlewares"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type invitesHandler struct {
	router *gin.Engine
	ctx    *repositories.Context
}

func (h *invitesHandler) RegisterRouter() {
	invite := h.router.Group("/invites")
	invite.Use(middlewares.AdminMiddleware(h.ctx.AdminRepository))
	invite.GET("/", h.get)
	invite.POST("/", h.post)
	invite.DELETE("/:id", h.deleteById)
	invite.PUT("/:id", h.putById)
}

func New(router *gin.Engine, ctx *repositories.Context) *invitesHandler {
	return &invitesHandler{
		router: router,
		ctx:    ctx,
	}
}
