package AuthHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type authHandler struct {
	router *gin.Engine
	ctx    *repositories.Context
}

func (h *authHandler) RegisterRouter() {
	auth := h.router.Group("/auth")
	auth.POST("/signup", h.postSignUp)
	auth.POST("/login", h.postLogin)
}

func New(router *gin.Engine, ctx *repositories.Context) *authHandler {
	return &authHandler{
		router: router,
		ctx:    ctx,
	}
}
