package tests

import (
	"github.com/stretchr/testify/suite"
	"github.com/the_fourth_dimension/planet_registry/pkg/database"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/routes"
)

type MiddlewareTestSuite struct {
	suite.Suite
	router *routes.Router
	db     *database.Database
}

// func (suite *Middlewar)
func makeAuthHeader(token string) string {
	return "Bearer " + token
}

func (suite *MiddlewareTestSuite) SetupTest() {
	suite.db.MigrateModels()
	suite.db.PopulateConfig()
}

func (suite *MiddlewareTestSuite) TearDownTest() {
	suite.db.DB.DropTable(&models.Admin{}, &models.Config{}, &models.Invite{}, &models.Planet{})
}
