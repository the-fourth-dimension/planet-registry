package tests

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

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

func (suite *SuperuserMiddlewareTestSuite) TestSuperuserMiddlewareWithAnInferiorRoleValue() {
	token, err := lib.SignJwt(j.MapClaims{"role": 2})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/superuser", nil)
	req.Header.Set("Authorization", lib.MakeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 401, w.Code)
}

type SuperuserMiddlewareTestSuite struct {
	suite.Suite
	router *routes.Router
	db     *database.Database
}

func (suite *SuperuserMiddlewareTestSuite) SetupTest() {
	suite.db.MigrateModels()
	suite.db.PopulateConfig()
}

func (suite *SuperuserMiddlewareTestSuite) TearDownTest() {
	suite.db.DB.DropTable(&models.Admin{}, &models.Config{}, &models.Invite{}, &models.Planet{})
}

func TestSuperuserMiddlewares(t *testing.T) {
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
	suite.Run(t, &SuperuserMiddlewareTestSuite{db: db, router: router})
	db.DB.Close()
}
