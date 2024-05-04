package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	AdminHandler "github.com/the_fourth_dimension/planet_registry/pkg/handlers/admin"
	InviteHandler "github.com/the_fourth_dimension/planet_registry/pkg/handlers/invite"
	PlanetHandler "github.com/the_fourth_dimension/planet_registry/pkg/handlers/planet"
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
	planetHandler := PlanetHandler.New(r.Engine, ctx)
	planetHandler.RegisterRouter()
	adminHandler := AdminHandler.New(r.Engine, ctx)
	adminHandler.RegisterRouter()
	inviteHandler := InviteHandler.New(r.Engine, ctx)
	inviteHandler.RegisterRouter()
}
