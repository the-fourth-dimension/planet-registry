package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

func serializeBody(body interface{}) *bytes.Buffer {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(body)
	return &b
}

func (suite *HandlersTestSuite) TestPlanetsPostHandlerWithPreExistingUsername() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	suite.ctx.PlanetRepository.Save(&models.Planet{PlanetId: "earth", Password: "password"})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password"})
	req, _ := http.NewRequest("POST", "/planets/", body)
	req.Header.Set("Authorization", makeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 409, w.Code)
}

func (suite *HandlersTestSuite) TestPlanetsPostHandlerWithValidInviteAndInviteOnlySetToTrue() {
	suite.ctx.ConfigRepository.Save(&models.Config{Model: gorm.Model{ID: 1}, InviteOnly: true})
	suite.ctx.InviteRepository.Save(&models.Invite{Code: "welcome"})
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password", "code": "welcome"})
	req, _ := http.NewRequest("POST", "/planets/", body)
	req.Header.Set("Authorization", makeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 201, w.Code)
	assert.Equal(suite.T(), len(suite.ctx.InviteRepository.Find(&models.Invite{}).Result), 0)
	assert.Nil(suite.T(), suite.ctx.PlanetRepository.FindFirst(&models.Planet{PlanetId: "earth"}).Error)
}

func (suite *HandlersTestSuite) TestPlanetsPostHandlerWithInvalidInviteAndInviteOnlySetToTrue() {
	suite.ctx.ConfigRepository.Save(&models.Config{Model: gorm.Model{ID: 1}, InviteOnly: true})
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password", "code": "welcome"})
	req, _ := http.NewRequest("POST", "/planets/", body)
	req.Header.Set("Authorization", makeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 400, w.Code)
	assert.NotNil(suite.T(), suite.ctx.PlanetRepository.FindFirst(&models.Planet{PlanetId: "earth"}).Error)
}

func (suite *HandlersTestSuite) TestPlanetsPostHandlerWithInvalidInviteAndInviteOnlySetToFalse() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password", "code": "welcome"})
	req, _ := http.NewRequest("POST", "/planets/", body)
	req.Header.Set("Authorization", makeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 201, w.Code)
	assert.Nil(suite.T(), suite.ctx.PlanetRepository.FindFirst(&models.Planet{PlanetId: "earth"}).Error)
}

func (suite *HandlersTestSuite) TestPlanetsPostLoginHandlerWithInvalidUsername() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password"})
	req, _ := http.NewRequest("POST", "/planets/login", body)
	req.Header.Set("Authorization", makeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *HandlersTestSuite) TestPlanetsPostLoginHandlerWithValidPassword() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	suite.ctx.PlanetRepository.Save(&models.Planet{PlanetId: "earth", Password: "password"})
	body := serializeBody(gin.H{"planetId": "earth", "password": "password"})
	req, _ := http.NewRequest("POST", "/planets/login", body)
	req.Header.Set("Authorization", makeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 200, w.Code)
}

func (suite *HandlersTestSuite) TestPlanetsPostLoginHandlerWithInvalidPassword() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	suite.ctx.PlanetRepository.Save(&models.Planet{PlanetId: "earth", Password: "password"})
	body := serializeBody(gin.H{"planetId": "earth", "password": "notpassword"})
	req, _ := http.NewRequest("POST", "/planets/login", body)
	req.Header.Set("Authorization", makeAuthHeader(token))

	suite.router.Engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), 403, w.Code)
}