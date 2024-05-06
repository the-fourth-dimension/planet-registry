package tests

import (
	"log"
	"net/http"
	"net/http/httptest"

	j "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
)

func (suite *MiddlewareTestSuite) TestSuperuserMiddlewareWithAnInferiorRoleValue() {
	token, err := lib.SignJwt(j.MapClaims{"role": 2})
	if err != nil {
		log.Panic("an error occurred while signing the token ", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/superuser", nil)
	req.Header.Set("Authorization", makeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 401, w.Code)
}
