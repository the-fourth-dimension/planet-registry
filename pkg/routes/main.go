package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/database"
)

type Router struct {
	Engine *gin.Engine
	DB     *database.Database
}

func NewRouter(db *database.Database) *Router {
	return &Router{
		Engine: gin.Default(),
		DB:     db,
	}
}
