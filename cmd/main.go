package main

import (
	"fmt"
	"log"

	"github.com/the_fourth_dimension/planet_registry/pkg/database"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
	"github.com/the_fourth_dimension/planet_registry/pkg/routes"
)

func main() {
	env.LoadEnv()
	logger := log.Default()
	db := database.ConnectDatabase(logger)
	db.MigrateModels()
	db.PopulateConfig()
	ctx := repositories.NewContext(db.DB, logger)

	router := routes.NewRouter(ctx)

	router.RegisterMiddlewares()
	router.RegisterRoutes()

	router.Engine.Run(fmt.Sprintf(":%s", env.GetEnv(env.PORT)))
}
