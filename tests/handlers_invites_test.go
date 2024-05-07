package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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

func (suite *InvitesHandlersTestSuite) TestInvitesGetWithNoCodeQuery() {
	w := httptest.NewRecorder()
	jsonFile, _ := os.Open("./mock_data/invites.json")
	byteValue, _ := io.ReadAll(jsonFile)
	var invites []models.Invite
	json.Unmarshal(byteValue, &invites)
	for _, invite := range invites {
		suite.ctx.InviteRepository.Save(&invite)
	}
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("GET", "/invites/", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	bodyBytes, _ := io.ReadAll(w.Result().Body)
	var receivedInvites []models.Invite
	json.Unmarshal(bodyBytes, &receivedInvites)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), len(invites), len(receivedInvites))
}

func (suite *InvitesHandlersTestSuite) TestInvitesGetWithNonExistingCodeQuery() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("GET", "/invites/?code=noob", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	bodyBytes, _ := io.ReadAll(w.Result().Body)
	var receivedInvites []models.Invite
	json.Unmarshal(bodyBytes, &receivedInvites)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), 0, len(receivedInvites))
}

func (suite *InvitesHandlersTestSuite) TestInvitesGetWithExistingCodeQuery() {
	w := httptest.NewRecorder()
	jsonFile, _ := os.Open("./mock_data/invites.json")
	byteValue, _ := io.ReadAll(jsonFile)
	var invites []models.Invite
	json.Unmarshal(byteValue, &invites)
	for _, invite := range invites {
		suite.ctx.InviteRepository.Save(&invite)
	}
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("GET", "/invites/?code=welcome", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	bodyBytes, _ := io.ReadAll(w.Result().Body)
	var receivedInvites []models.Invite
	json.Unmarshal(bodyBytes, &receivedInvites)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), 1, len(receivedInvites))
}

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
