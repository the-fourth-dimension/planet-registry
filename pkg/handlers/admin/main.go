package AdminHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/middlewares"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type credentials struct {
	username string
	password string
}

type adminHandler struct {
	router *gin.Engine
	ctx    *repositories.Context
}

func (h *adminHandler) RegisterRouter() {
	admin := h.router.Group("/admin")
	admin.Use(middlewares.SuperuserMiddleware())
	admin.GET("/", h.get)
	admin.POST("/", h.post)
	admin.PUT("/:id", h.putById)
	admin.DELETE("/:id", h.deleteById)
}

func New(router *gin.Engine, ctx *repositories.Context) *adminHandler {
	return &adminHandler{
		router: router,
		ctx:    ctx,
	}
}
