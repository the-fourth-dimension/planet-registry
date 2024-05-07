package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/the_fourth_dimension/planet_registry/pkg/lib"
	"github.com/the_fourth_dimension/planet_registry/pkg/models"
)

func (suite *HandlersTestSuite) TestAdminsGet() {
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
	req.Header.Set("Authorization", makeAuthHeader(token))
	suite.router.Engine.ServeHTTP(w, req)
	bodyBytes, _ := io.ReadAll(w.Result().Body)
	var receivedAdmins []models.Admin
	json.Unmarshal(bodyBytes, &receivedAdmins)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), len(admins), len(receivedAdmins))
}
