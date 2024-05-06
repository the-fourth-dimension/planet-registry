package tests

import (
	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/the_fourth_dimension/planet_registry/pkg/database"
	"github.com/the_fourth_dimension/planet_registry/pkg/env"
	"github.com/the_fourth_dimension/planet_registry/pkg/routes"
)

func TestHealthRoute(t *testing.T) {
	os.Setenv("APP_ENV", "TEST")
	env.LoadEnv()
	db := database.ConnectDatabase()
	router := routes.NewRouter(db.DB)
	router.RegisterMiddlewares()
	router.RegisterRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
