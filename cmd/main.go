package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
	"github.com/the_fourth_dimension/planet_registry/pkg/middlewares"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/routes"
)

func main() {
	env.LoadEnv()
	models.ConnectDatabase()

	router := gin.Default()

	router.Use(cors.Default())
	router.Use(middlewares.LoggerMiddleware())

	routes.IntializeRoutes(router)

	router.Run(fmt.Sprintf(":%s", env.GetEnv(env.PORT)))
}
