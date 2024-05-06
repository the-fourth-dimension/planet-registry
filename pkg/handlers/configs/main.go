package ConfigsHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/middlewares"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type configsHandler struct {
	router *gin.Engine
	ctx    *repositories.Context
}

func (h *configsHandler) RegisterRouter() {
	admin := h.router.Group("/configs")
	admin.Use(middlewares.SuperuserMiddleware())
	admin.GET("/", h.get)
}

func New(router *gin.Engine, ctx *repositories.Context) *configsHandler {
	return &configsHandler{
		router: router,
		ctx:    ctx,
	}
}
