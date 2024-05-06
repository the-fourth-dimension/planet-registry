package ConfigHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/middlewares"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type configHandler struct {
	router *gin.Engine
	ctx    *repositories.Context
}

func (h *configHandler) RegisterRouter() {
	admin := h.router.Group("/configs")
	admin.Use(middlewares.SuperuserMiddleware())
}

func New(router *gin.Engine, ctx *repositories.Context) *configHandler {
	return &configHandler{
		router: router,
		ctx:    ctx,
	}
}
