package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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

func serializeBody(body interface{}) *bytes.Buffer {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(body)
	return &b
}

func (suite *PlanetsHandlersTestSuite) TestPlanetsPostHandlerWithPreExistingUsername() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	suite.ctx.PlanetRepository.Save(&models.Planet{PlanetId: "earth", Password: "password"})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password"})
	req, _ := http.NewRequest("POST", "/planets/", body)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 409, w.Code)
}

func (suite *PlanetsHandlersTestSuite) TestPlanetsPostHandlerWithValidInviteAndInviteOnlySetToTrue() {
	suite.ctx.ConfigRepository.Save(&models.Config{Model: gorm.Model{ID: 1}, InviteOnly: true})
	suite.ctx.InviteRepository.Save(&models.Invite{Code: "welcome"})
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password", "code": "welcome"})
	req, _ := http.NewRequest("POST", "/planets/", body)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 201, w.Code)
	assert.Equal(suite.T(), len(suite.ctx.InviteRepository.Find(&models.Invite{}).Result), 0)
	assert.Nil(suite.T(), suite.ctx.PlanetRepository.FindFirst(&models.Planet{PlanetId: "earth"}).Error)
}

func (suite *PlanetsHandlersTestSuite) TestPlanetsPostHandlerWithInvalidInviteAndInviteOnlySetToTrue() {
	suite.ctx.ConfigRepository.Save(&models.Config{Model: gorm.Model{ID: 1}, InviteOnly: true})
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password", "code": "welcome"})
	req, _ := http.NewRequest("POST", "/planets/", body)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 400, w.Code)
	assert.NotNil(suite.T(), suite.ctx.PlanetRepository.FindFirst(&models.Planet{PlanetId: "earth"}).Error)
}

func (suite *PlanetsHandlersTestSuite) TestPlanetsPostHandlerWithInvalidInviteAndInviteOnlySetToFalse() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password", "code": "welcome"})
	req, _ := http.NewRequest("POST", "/planets/", body)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 201, w.Code)
	assert.Nil(suite.T(), suite.ctx.PlanetRepository.FindFirst(&models.Planet{PlanetId: "earth"}).Error)
}

func (suite *PlanetsHandlersTestSuite) TestPlanetsPostLoginHandlerWithInvalidUsername() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password"})
	req, _ := http.NewRequest("POST", "/planets/login", body)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *PlanetsHandlersTestSuite) TestPlanetsPostLoginHandlerWithValidPassword() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	suite.ctx.PlanetRepository.Save(&models.Planet{PlanetId: "earth", Password: "password"})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password"})
	req, _ := http.NewRequest("POST", "/planets/login", body)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 200, w.Code)
}

func (suite *PlanetsHandlersTestSuite) TestPlanetsPostLoginHandlerWithInvalidPassword() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	suite.ctx.PlanetRepository.Save(&models.Planet{PlanetId: "earth", Password: "password"})
	body := serializeBody(gin.H{"planetId": "earth", "password": "notpassword"})
	req, _ := http.NewRequest("POST", "/planets/login", body)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 403, w.Code)
}

func (suite *PlanetsHandlersTestSuite) SetupTest() {
	suite.db.MigrateModels()
	suite.db.PopulateConfig()
	suite.ctx.ConfigRepository.Save(&models.Config{Model: gorm.Model{ID: 1}, InviteOnly: false})
}

func (suite *PlanetsHandlersTestSuite) TearDownTest() {
	suite.db.DB.DropTable(&models.Admin{}, &models.Config{}, &models.Invite{}, &models.Planet{})
}

type PlanetsHandlersTestSuite struct {
	suite.Suite
	router *routes.Router
	db     *database.Database
	ctx    *repositories.Context
}

func TestPlanetsHandlers(t *testing.T) {
	logger := log.Default()
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	env.LoadEnv()
	db := database.ConnectDatabase(logger)
	ctx := repositories.NewContext(db.DB, logger)
	router := routes.NewRouter(ctx)
	router.RegisterMiddlewares()
	router.RegisterRoutes()
	suite.Run(t, &PlanetsHandlersTestSuite{db: db, router: router, ctx: ctx})
	db.DB.Close()
}
