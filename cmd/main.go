package main

import (
	"fmt"

	"github.com/the_fourth_dimension/planet_registry/pkg/database"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
	"github.com/the_fourth_dimension/planet_registry/pkg/routes"
)

func main() {
	env.LoadEnv()
	db := database.ConnectDatabase()
	db.MigrateModels()
	db.PopulateConfig()
	router := routes.NewRouter(db)

	router.RegisterMiddlewares()
	router.RegisterRoutes()

	router.Engine.Run(fmt.Sprintf(":%s", env.GetEnv(env.PORT)))
}
