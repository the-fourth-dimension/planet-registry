package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	AdminsHandler "github.com/the_fourth_dimension/planet_registry/pkg/handlers/admins"
	ConfigsHandler "github.com/the_fourth_dimension/planet_registry/pkg/handlers/config"
	InvitesHandler "github.com/the_fourth_dimension/planet_registry/pkg/handlers/invites"
	PlanetsHandler "github.com/the_fourth_dimension/planet_registry/pkg/handlers/planet"
	"github.com/the_fourth_dimension/planet_registry/pkg/middlewares"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
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
	ctx := repositories.NewContext(r.DB)
	planetHandler := PlanetsHandler.New(r.Engine, ctx)
	planetHandler.RegisterRouter()
	adminHandler := AdminsHandler.New(r.Engine, ctx)
	adminHandler.RegisterRouter()
	inviteHandler := InvitesHandler.New(r.Engine, ctx)
	inviteHandler.RegisterRouter()
	configHandler := ConfigsHandler.New(r.Engine, ctx)
	configHandler.RegisterRouter()
}
