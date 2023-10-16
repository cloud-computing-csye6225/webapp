package integration_tests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapp/config"
	"webapp/db"
	"webapp/routes"
)

func TestHealthzRouteWithGET(t *testing.T) {
	configs := config.GetConfigs()
	pgDB := &db.PostgresDB{}
	err := pgDB.InitDatabase(configs.DBConfig)
	if err != nil {
		fmt.Printf("failed to initialize database: %s", err)
	}

	router := gin.Default()
	router.Any("/healthz", routes.HealthzGetReqHandler(pgDB))

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Got %d want %d", response.Code, http.StatusOK)
	}
}
