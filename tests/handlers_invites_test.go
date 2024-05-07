package tests

import (
	"bytes"
	"log"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	"github.com/the_fourth_dimension/planet_registry/pkg/database"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
	"github.com/the_fourth_dimension/planet_registry/pkg/routes"
)

func (suite *InvitesHandlersTestSuite) SetupTest() {
	suite.db.MigrateModels()
	suite.db.PopulateConfig()
	suite.ctx.ConfigRepository.Save(&models.Config{Model: gorm.Model{ID: 1}, InviteOnly: false})
}

func (suite *InvitesHandlersTestSuite) TearDownTest() {
	suite.db.DB.DropTable(&models.Admin{}, &models.Config{}, &models.Invite{}, &models.Planet{})
}

type InvitesHandlersTestSuite struct {
	suite.Suite
	router *routes.Router
	db     *database.Database
	ctx    *repositories.Context
}

func TestInvitesHandlers(t *testing.T) {
	logger := log.Default()
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	env.LoadEnv()
	db := database.ConnectDatabase(logger)
	ctx := repositories.NewContext(db.DB, logger)
	router := routes.NewRouter(ctx)
	router.RegisterMiddlewares()
	router.RegisterRoutes()
	suite.Run(t, &InvitesHandlersTestSuite{db: db, router: router, ctx: ctx})
	db.DB.Close()
}
