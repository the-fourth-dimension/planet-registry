package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
)

type Router struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

var appEnvToGinMode = map[string]string{
	"TEST":       gin.TestMode,
	"DEV":        gin.DebugMode,
	"PRODUCTION": gin.ReleaseMode,
}

func NewRouter(db *gorm.DB) *Router {
	gin.SetMode(appEnvToGinMode[env.GetEnv(env.APP_ENV)])
	return &Router{
		Engine: gin.Default(),
		DB:     db,
	}
}
