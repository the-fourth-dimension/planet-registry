package AdminsHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/middlewares"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type adminsHandler struct {
	router *gin.Engine
	ctx    *repositories.Context
}

func (h *adminsHandler) RegisterRouter() {
	admin := h.router.Group("/admins")
	admin.Use(middlewares.SuperuserMiddleware())
	admin.GET("/", h.get)
	admin.POST("/", h.post)
	admin.PUT("/:id", h.putById)
	admin.DELETE("/:id", h.deleteById)
}

func New(router *gin.Engine, ctx *repositories.Context) *adminsHandler {
	return &adminsHandler{
		router: router,
		ctx:    ctx,
	}
}
