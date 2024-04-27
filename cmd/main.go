package main

import (
	"fmt"

	"github.com/the_fourth_dimension/planet_registry/pkg/env"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/routes"
)

func main() {
	env.LoadEnv()
	models.ConnectDatabase()

	router := routes.NewRouter()

	router.RegisterMiddlewares()
	router.RegisterRoutes()

	router.Engine.Run(fmt.Sprintf(":%s", env.GetEnv(env.PORT)))
}
