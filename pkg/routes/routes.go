package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	AdminsHandler "github.com/the_fourth_dimension/planet_registry/pkg/handlers/admins"
	ConfigsHandler "github.com/the_fourth_dimension/planet_registry/pkg/handlers/configs"
	InvitesHandler "github.com/the_fourth_dimension/planet_registry/pkg/handlers/invites"
	PlanetsHandler "github.com/the_fourth_dimension/planet_registry/pkg/handlers/planet"
	"github.com/the_fourth_dimension/planet_registry/pkg/middlewares"
)

func (r *Router) RegisterMiddlewares() {
	r.Engine.Use(cors.Default())
	r.Engine.Use(middlewares.ErrorMiddleware())
	r.Engine.Use(middlewares.LoggerMiddleware())
	r.Engine.Use(middlewares.AuthMiddleware())
}

func (r *Router) RegisterRoutes() {

	r.Engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	r.Engine.GET("/", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	planetHandler := PlanetsHandler.New(r.Engine, r.ctx)
	planetHandler.RegisterRouter()
	adminHandler := AdminsHandler.New(r.Engine, r.ctx)
	adminHandler.RegisterRouter()
	inviteHandler := InvitesHandler.New(r.Engine, r.ctx)
	inviteHandler.RegisterRouter()
	configHandler := ConfigsHandler.New(r.Engine, r.ctx)
	configHandler.RegisterRouter()
}
