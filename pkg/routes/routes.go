package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/handlers"
	"github.com/the_fourth_dimension/planet_registry/pkg/middlewares"
)

func (r *Router) RegisterMiddlewares() {
	r.Engine.Use(cors.Default())
	r.Engine.Use(middlewares.LoggerMiddleware())
}

func (r *Router) RegisterRoutes() {

	r.Engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	r.Engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Hi")
	})

	r.Engine.POST("/auth/signup", handlers.Signup)

	r.Engine.POST("/auth/login", handlers.Login)

}
