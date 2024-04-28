package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	handlers_auth "github.com/the_fourth_dimension/planet_registry/pkg/handlers/auth"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

func NewAuthHandler(router *gin.Engine, db *gorm.DB) *handlers_auth.AuthHandler {
	return &handlers_auth.AuthHandler{
		PlanetRepository:     *repositories.NewPlanetRepository(db),
		InviteCodeRepository: *repositories.NewInviteCodeRepository(db),
		ConfigRepository:     *repositories.NewConfigRepository(db),
		Router:               router,
	}
}
