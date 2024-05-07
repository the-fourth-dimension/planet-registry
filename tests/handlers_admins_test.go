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

func (suite *AdminsHandlersTestSuite) TestAdminsGet() {
	w := httptest.NewRecorder()
	jsonFile, _ := os.Open("./mock_data/admins.json")
	byteValue, _ := io.ReadAll(jsonFile)
	var admins []models.Admin
	json.Unmarshal(byteValue, &admins)
	for _, admin := range admins {
		suite.ctx.AdminRepository.Save(&admin)
	}
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("GET", "/admins/", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	bodyBytes, _ := io.ReadAll(w.Result().Body)
	var receivedAdmins []models.Admin
	json.Unmarshal(bodyBytes, &receivedAdmins)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), len(admins), len(receivedAdmins))
}

func (suite *AdminsHandlersTestSuite) TestAdminsPostWithPreExistingUsername() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	body := lib.SerializeBody(gin.H{"username": "probablyarth", "password": "password"})
	req, _ := http.NewRequest("POST", "/admins/", body)
	suite.ctx.AdminRepository.Save(&models.Admin{Username: "probablyarth", Password: "Password"})
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 409, w.Code)
}

func (suite *AdminsHandlersTestSuite) TestAdminsPostWithValidUsername() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	body := lib.SerializeBody(gin.H{"username": "probablyarth", "password": "password"})
	req, _ := http.NewRequest("POST", "/admins/", body)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 201, w.Code)
}

func (suite *AdminsHandlersTestSuite) TestAdminsPutIdWithNoIdParam() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("PUT", "/admins/", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *AdminsHandlersTestSuite) TestAdminsPutIdWithIdParamOfNonIntegerDatatype() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("PUT", "/admins/arth", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *AdminsHandlersTestSuite) TestAdminsPutIdWithNonExistingIdParam() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("PUT", "/admins/1", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *AdminsHandlersTestSuite) TestAdminsPutIdWithEmptyJsonBody() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	suite.ctx.AdminRepository.Save(&models.Admin{Username: "probablyarth", Password: "password"})
	req, _ := http.NewRequest("PUT", "/admins/1", lib.SerializeBody(gin.H{}))
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 204, w.Code)
}

func (suite *AdminsHandlersTestSuite) TestAdminsPutIdWithDifferentPassword() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	suite.ctx.AdminRepository.Save(&models.Admin{Username: "probablyarth", Password: "password"})
	newPassword := "new_password"
	req, _ := http.NewRequest("PUT", "/admins/1", lib.SerializeBody(gin.H{"password": newPassword}))
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	isMatch := lib.VerifyPassword(newPassword, suite.ctx.AdminRepository.FindFirst(&models.Admin{Model: gorm.Model{ID: 1}}).Result.Password)
	assert.True(suite.T(), isMatch == nil)
	assert.Equal(suite.T(), 202, w.Code)
}

func (suite *AdminsHandlersTestSuite) TestAdminsPutIdWithDifferentUsername() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	suite.ctx.AdminRepository.Save(&models.Admin{Username: "probablyarth", Password: "password"})
	newUsername := "definitelyarth"
	req, _ := http.NewRequest("PUT", "/admins/1", lib.SerializeBody(gin.H{"username": newUsername}))
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), newUsername, suite.ctx.AdminRepository.FindFirst(&models.Admin{Model: gorm.Model{ID: 1}}).Result.Username)
	assert.Equal(suite.T(), 202, w.Code)
}

func (suite *AdminsHandlersTestSuite) TestAdminsDeleteIdWithNoIdParam() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("DELETE", "/admins/", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *AdminsHandlersTestSuite) TestAdminsDeleteIdWithIdParamOfNonIntegerDatatype() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("DELETE", "/admins/arth", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *AdminsHandlersTestSuite) TestAdminsDeleteIdWithNonExistingIdParam() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("DELETE", "/admins/1", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *AdminsHandlersTestSuite) TestAdminsDeleteIdWithValidIdParam() {
	w := httptest.NewRecorder()

	token, _ := lib.SignJwt(jwt.MapClaims{"role": 0})
	req, _ := http.NewRequest("DELETE", "/admins/1", nil)
	suite.ctx.AdminRepository.Save(&models.Admin{Username: "probablyarth", Password: "password"})
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)

	assert.True(suite.T(), suite.ctx.AdminRepository.FindFirst(&models.Admin{Model: gorm.Model{ID: 1}}).Error != nil)
	assert.Equal(suite.T(), 204, w.Code)
}

func (suite *AdminsHandlersTestSuite) SetupTest() {
	suite.db.MigrateModels()
	suite.db.PopulateConfig()
	suite.ctx.ConfigRepository.Save(&models.Config{Model: gorm.Model{ID: 1}, InviteOnly: false})
}

func (suite *AdminsHandlersTestSuite) TearDownTest() {
	suite.db.DB.DropTable(&models.Admin{}, &models.Config{}, &models.Invite{}, &models.Planet{})
}

type AdminsHandlersTestSuite struct {
	suite.Suite
	router *routes.Router
	db     *database.Database
	ctx    *repositories.Context
}

func TestAdminHandlers(t *testing.T) {
	logger := log.Default()
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	env.LoadEnv()
	db := database.ConnectDatabase(logger)
	ctx := repositories.NewContext(db.DB, logger)
	router := routes.NewRouter(ctx)
	router.RegisterMiddlewares()
	router.RegisterRoutes()
	suite.Run(t, &AdminsHandlersTestSuite{db: db, router: router, ctx: ctx})
	db.DB.Close()
}
