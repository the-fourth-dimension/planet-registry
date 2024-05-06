package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

type Router struct {
	Engine *gin.Engine
	ctx    *repositories.Context
}

var appEnvToGinMode = map[string]string{
	"TEST":       gin.TestMode,
	"DEV":        gin.DebugMode,
	"PRODUCTION": gin.ReleaseMode,
}

func NewRouter(ctx *repositories.Context) *Router {
	gin.SetMode(appEnvToGinMode[env.GetEnv(env.APP_ENV)])
	return &Router{
		Engine: gin.New(),
		ctx:    ctx,
	}
}
