package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/handlers"
)

func IntializeRoutes(r *gin.Engine) {

	r.GET("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Hi")
	})

	r.POST("/auth/signup", handlers.Signup)

	r.POST("/auth/login", handlers.Login)

}
