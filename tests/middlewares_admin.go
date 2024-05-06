package tests

import (
	"log"
	"net/http"
	"net/http/httptest"

	j "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
	"github.com/the_fourth_dimension/planet_registry/pkg/repositories"
)

func (suite *MiddlewareTestSuite) TestAdminMiddlewareWithAnInferiorRoleValue() {
	token, err := lib.SignJwt(j.MapClaims{"role": 2, "username": "probablyarth"})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", makeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 401, w.Code)
}

func (suite *MiddlewareTestSuite) TestAdminMiddlewareWithValidRoleValueAndInvalidUsername() {
	token, err := lib.SignJwt(j.MapClaims{"role": 1, "username": "probablyarth"})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", makeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 403, w.Code)
}

func (suite *MiddlewareTestSuite) TestAdminMiddlewareWithValidRoleValueAndValidUsername() {
	adminRepository := repositories.NewAdminRepository(suite.db.DB)
	saveResult := adminRepository.Save(&models.Admin{Username: "probablyarth", Password: "password"})
	if saveResult.Error != nil {
		panic(saveResult.Error)
	}
	token, err := lib.SignJwt(j.MapClaims{"role": 1, "username": "probablyarth"})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", makeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 200, w.Code)
}
