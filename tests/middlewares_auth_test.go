package tests

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	j "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/the_fourth_dimension/planet_registry/pkg/database"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/middlewares"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
	"github.com/the_fourth_dimension/planet_registry/pkg/routes"
)

func testHandler(ctx *gin.Context) {
	ctx.Status(200)
}
func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareWithMissingToken() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareWithMissingRoleClaim() {
	token, err := lib.SignJwt(j.MapClaims{})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareWithInvalidRoleClaim() {
	token, err := lib.SignJwt(j.MapClaims{"role": 10})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 401, w.Code)
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareWithValidRoleClaim() {
	token, err := lib.SignJwt(j.MapClaims{"role": 0})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 200, w.Code)
}

func (suite *AuthMiddlewareTestSuite) TestAuthMiddlewareWithAdminRoleAndMissingUsernameClaim() {
	token, err := lib.SignJwt(j.MapClaims{"role": 1})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 403, w.Code)
}

type AuthMiddlewareTestSuite struct {
	suite.Suite
	router *routes.Router
	db     *database.Database
}

func (suite *AuthMiddlewareTestSuite) SetupTest() {
	suite.db.MigrateModels()
	suite.db.PopulateConfig()
}

func (suite *AuthMiddlewareTestSuite) TearDownTest() {
	suite.db.DB.DropTable(&models.Admin{}, &models.Config{}, &models.Invite{}, &models.Planet{})
}

func TestAuthMiddlewares(t *testing.T) {
	logger := log.Default()
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	env.LoadEnv()
	db := database.ConnectDatabase(logger)
	ctx := repositories.NewContext(db.DB, logger)
	router := routes.NewRouter(ctx)
	router.Engine.GET("/auth", middlewares.ErrorMiddleware(), middlewares.AuthMiddleware(), testHandler)
	router.Engine.GET("/admin", middlewares.ErrorMiddleware(), middlewares.AuthMiddleware(), middlewares.AdminMiddleware(ctx.AdminRepository), testHandler)
	router.Engine.GET("/superuser", middlewares.ErrorMiddleware(), middlewares.AuthMiddleware(), middlewares.SuperuserMiddleware(), testHandler)
	suite.Run(t, &AuthMiddlewareTestSuite{db: db, router: router})
	db.DB.Close()
}
