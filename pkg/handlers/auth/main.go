package handlers_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type Credentials struct {
	PlanetId string `json:"planetId" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CredentialsWithCode struct {
	Credentials
	Code string `json:"code" binding:"required"`
}

type AuthHandler struct {
	PlanetRepository     repositories.PlanetRepository
	InviteCodeRepository repositories.InviteCodeRepository
	ConfigRepository     repositories.ConfigRepository
	Router               *gin.Engine
}

// [ ] USE Transaction in signup
// [x] separate methods for models into repositories
func (h *AuthHandler) RegisterRouter() {
	auth := h.Router.Group("/auth")
	auth.POST("/signup", h.postSignUp)
	auth.POST("/login", h.postLogin)
}
