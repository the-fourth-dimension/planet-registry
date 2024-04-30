package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/database"
	handlers_auth "github.com/the_fourth_dimension/planet_registry/pkg/handlers/auth"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

func NewAuthHandler(router *gin.Engine, db *database.Database) *handlers_auth.AuthHandler {
	return &handlers_auth.AuthHandler{
		PlanetRepository:     *repositories.NewPlanetRepository(db.DB),
		InviteCodeRepository: *repositories.NewInviteCodeRepository(db.DB),
		ConfigRepository:     *repositories.NewConfigRepository(db.DB),
		Router:               router,
		ExecuteTransaction:   db.ExecuteTransaction,
	}
}
