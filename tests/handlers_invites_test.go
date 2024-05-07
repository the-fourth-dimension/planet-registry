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

func (suite *InvitesHandlersTestSuite) TestInvitesPostWithPreExistingCode() {
	w := httptest.NewRecorder()
	jsonFile, _ := os.Open("./mock_data/invites.json")
	byteValue, _ := io.ReadAll(jsonFile)
	var invites []models.Invite
	json.Unmarshal(byteValue, &invites)
	for _, invite := range invites {
		suite.ctx.InviteRepository.Save(&invite)
	}
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("POST", "/invites/", lib.SerializeBody(gin.H{"code": "welcome"}))
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 409, w.Code)
}

func (suite *InvitesHandlersTestSuite) TestInvitesPostWithInvalidCodeLength() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("POST", "/invites/", lib.SerializeBody(gin.H{"code": "wel"}))
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *InvitesHandlersTestSuite) TestInvitesPostWithValidCodeLength() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("POST", "/invites/", lib.SerializeBody(gin.H{"code": "welcome"}))
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	bodyBytes, _ := io.ReadAll(w.Result().Body)
	var receivedInvite struct {
		Id int `json:"id"`
	}
	json.Unmarshal(bodyBytes, &receivedInvite)
	assert.Equal(suite.T(), 201, w.Code)
	assert.Equal(suite.T(), 1, receivedInvite.Id)
	assert.True(suite.T(), suite.ctx.InviteRepository.FindFirst(&models.Invite{Code: "welcome"}).Error == nil)
}

func (suite *InvitesHandlersTestSuite) TestInvitesDeleteIdWithNoIdParam() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("DELETE", "/invites/", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *InvitesHandlersTestSuite) TestInvitesDeleteIdWithIdParamOfNonIntegerDatatype() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("DELETE", "/invites/arth", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *InvitesHandlersTestSuite) TestInvitesDeleteIdWithNonExistingIdParam() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("DELETE", "/invites/1", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *InvitesHandlersTestSuite) TestInvitesDeleteIdWithValidIdParam() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("DELETE", "/invites/1", nil)
	suite.ctx.InviteRepository.Save(&models.Invite{Code: "welcome"})
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.True(suite.T(), suite.ctx.InviteRepository.FindFirst(&models.Invite{Model: gorm.Model{ID: 1}}).Error != nil)
	assert.Equal(suite.T(), 204, w.Code)
}

func (suite *InvitesHandlersTestSuite) SetupTest() {
	suite.db.MigrateModels()
	suite.db.PopulateConfig()
	suite.ctx.ConfigRepository.Save(&models.Config{Model: gorm.Model{ID: 1}, InviteOnly: false})
}

func (suite *InvitesHandlersTestSuite) TearDownTest() {
	suite.db.DB.DropTable(&models.Admin{}, &models.Config{}, &models.Invite{}, &models.Planet{})
}

func (suite *InvitesHandlersTestSuite) TestInvitesPutIdWithNoIdParam() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("PUT", "/invites/", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *InvitesHandlersTestSuite) TestInvitesPutIdWithIdParamOfNonIntegerDatatype() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("PUT", "/invites/arth", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *InvitesHandlersTestSuite) TestInvitesPutIdWithNonExistingIdParam() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("PUT", "/invites/1", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *InvitesHandlersTestSuite) TestInvitesPutIdWithInvalidCodeLength() {
	w := httptest.NewRecorder()
	suite.ctx.InviteRepository.Save(&models.Invite{Code: "welcome"})
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("PUT", "/invites/1", lib.SerializeBody(gin.H{"code": "wel"}))
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 400, w.Code)
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
