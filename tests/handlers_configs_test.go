package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/the_fourth_dimension/planet_registry/pkg/database"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
	"github.com/the_fourth_dimension/planet_registry/pkg/routes"
)

func (suite *ConfigsHandlersTestSuite) TestConfigsGetWithNoConfigsRecord() {
	w := httptest.NewRecorder()
	suite.db.DB.Delete(&models.Config{Model: gorm.Model{ID: 1}})
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("GET", "/configs/", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	defer func() {
		r := recover()
		assert.NotNil(suite.T(), r)
	}()
	suite.router.Engine.ServeHTTP(w, req)
}

func (suite *ConfigsHandlersTestSuite) TestConfigsGetWithValidConfigsRecord() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("GET", "/configs/", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	bodyBytes, _ := io.ReadAll(w.Result().Body)
	var receivedConfig struct {
		Config models.Config `json:"config"`
	}
	json.Unmarshal(bodyBytes, &receivedConfig)
	assert.Equal(suite.T(), uint(1), receivedConfig.Config.ID)
}

func (suite *ConfigsHandlersTestSuite) SetupTest() {
	suite.db.MigrateModels()
	suite.db.PopulateConfig()
	suite.ctx.ConfigRepository.Save(&models.Config{Model: gorm.Model{ID: 1}, InviteOnly: false})
}

func (suite *ConfigsHandlersTestSuite) TearDownTest() {
	suite.db.DB.DropTable(&models.Admin{}, &models.Config{}, &models.Invite{}, &models.Planet{})
}

type ConfigsHandlersTestSuite struct {
	suite.Suite
	router *routes.Router
	db     *database.Database
	ctx    *repositories.Context
}

func TestConfigHandlers(t *testing.T) {
	logger := log.Default()
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	env.LoadEnv()
	db := database.ConnectDatabase(logger)
	ctx := repositories.NewContext(db.DB, logger)
	router := routes.NewRouter(ctx)
	router.RegisterMiddlewares()
	router.RegisterRoutes()
	suite.Run(t, &ConfigsHandlersTestSuite{db: db, router: router, ctx: ctx})
	db.DB.Close()
}
