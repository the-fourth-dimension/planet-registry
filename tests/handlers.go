package tests

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	"github.com/the_fourth_dimension/planet_registry/pkg/database"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
	"github.com/the_fourth_dimension/planet_registry/pkg/routes"
)

type HandlersTestSuite struct {
	suite.Suite
	router *routes.Router
	db     *database.Database
	ctx    *repositories.Context
}

func (suite *HandlersTestSuite) SetupTest() {
	suite.db.MigrateModels()
	suite.db.PopulateConfig()
	suite.ctx.ConfigRepository.Save(&models.Config{Model: gorm.Model{ID: 1}, InviteOnly: false})
}

func (suite *HandlersTestSuite) TearDownTest() {
	suite.db.DB.DropTable(&models.Admin{}, &models.Config{}, &models.Invite{}, &models.Planet{})
}
