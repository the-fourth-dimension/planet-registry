package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

func serializeBody(body interface{}) *bytes.Buffer {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(body)
	return &b
}

func (suite *HandlersTestSuite) TestPostPlanetsHandlerWithPreExistingUsername() {
	w := httptest.NewRecorder()
	token, _ := lib.SignJwt(jwt.MapClaims{"role": 2})
	suite.ctx.PlanetRepository.Save(&models.Planet{PlanetId: "earth", Password: "password"})
	body := serializeBody(map[string]string{"planetId": "earth", "password": "password"})
	req, _ := http.NewRequest("POST", "/planets/", body)
	req.Header.Set("Authorization", makeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	assert.Equal(suite.T(), 409, w.Code)
}
