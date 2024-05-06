package tests

import (
	"bytes"
	"log"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/the_fourth_dimension/planet_registry/pkg/database"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
	"github.com/the_fourth_dimension/planet_registry/pkg/middlewares"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
	"github.com/the_fourth_dimension/planet_registry/pkg/routes"
)

func testHandler(ctx *gin.Context) {
	ctx.Status(200)
}

func TestMiddlewares(t *testing.T) {
	logger := log.Default()
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	env.LoadEnv()
	db := database.ConnectDatabase(logger)
	ctx := repositories.NewContext(db.DB, logger)
	router := routes.NewRouter(ctx)
	router.Engine.GET("/auth", middlewares.ErrorMiddleware(), middlewares.AuthMiddleware(), testHandler)
	router.Engine.GET("/admin", middlewares.ErrorMiddleware(), middlewares.AuthMiddleware(), middlewares.AdminMiddleware(ctx.AdminRepository), testHandler)
	suite.Run(t, &MiddlewareTestSuite{db: db, router: router})
	db.DB.Close()
}
