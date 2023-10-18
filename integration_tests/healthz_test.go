package integration_tests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"webapp/config"
	"webapp/db"
	"webapp/routes"
)

func loadEnv() {
	const projectDirName = "webapp"
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		fmt.Printf("No .env file")
	}
}

func TestHealthzRouteWithGET(t *testing.T) {
	loadEnv()
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
