package PlanetHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type credentials struct {
	PlanetId string `json:"planetId" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type credentialsWithCode struct {
	credentials
	Code string `json:"code" binding:"required"`
}

type planetHandler struct {
	router *gin.Engine
	ctx    *repositories.Context
}

func (h *planetHandler) RegisterRouter() {
	auth := h.router.Group("/planet")
	auth.POST("/", h.post)
	auth.POST("/login", h.postLogin)
}

func New(router *gin.Engine, ctx *repositories.Context) *planetHandler {
	return &planetHandler{
		router: router,
		ctx:    ctx,
	}
}
