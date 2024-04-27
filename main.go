package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/the_fourth_dimension/planet_registry/middlewares"
	"github.com/the_fourth_dimension/planet_registry/models"
	"github.com/the_fourth_dimension/planet_registry/routes"
)

func main() {
	models.ConnectDatabase()

	router := gin.Default()
	
	// Middlewares
	router.Use(cors.Default())
	router.Use(middlewares.LoggerMiddleware())
	
	routes.IntializeRoutes(router)
	
	router.Run(":3000")
}