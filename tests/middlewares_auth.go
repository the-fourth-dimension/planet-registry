package tests

import (
	"log"
	"net/http"
	"net/http/httptest"

	j "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
)

func (suite *MiddlewareTestSuite) TestAuthMiddlewareWithMissingToken() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *MiddlewareTestSuite) TestAuthMiddlewareWithMissingRoleClaim() {
	token, err := lib.SignJwt(j.MapClaims{})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	req.Header.Set("Authorization", makeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *MiddlewareTestSuite) TestAuthMiddlewareWithInvalidRoleClaim() {
	token, err := lib.SignJwt(j.MapClaims{"role": 10})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	req.Header.Set("Authorization", makeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 401, w.Code)
}

func (suite *MiddlewareTestSuite) TestAuthMiddlewareWithValidRoleClaim() {
	token, err := lib.SignJwt(j.MapClaims{"role": 0})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	req.Header.Set("Authorization", makeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 200, w.Code)
}

func (suite *MiddlewareTestSuite) TestAuthMiddlewareWithAdminRoleAndMissingUsernameClaim() {
	token, err := lib.SignJwt(j.MapClaims{"role": 1})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	req.Header.Set("Authorization", makeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 403, w.Code)
}
