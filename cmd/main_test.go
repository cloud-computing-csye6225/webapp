package main

import (
	"errors"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"venkatpavan_munaganti_002722397_2/config"
)

type MockDB struct {
	Error       error
	isConnected bool
}

func (m *MockDB) InitDatabase(config.DatabaseConfig) error {
	if m.Error != nil {
		m.isConnected = false
		return m.Error
	}
	m.isConnected = true
	return m.Error
}

func (m *MockDB) GetConnection() *gorm.DB {
	return &gorm.DB{}
}

func (m *MockDB) Ping() error {
	if m.isConnected {
		return nil
	}
	return errors.New("database is not active")
}

func TestHealthzRouteWithGET(t *testing.T) {
	mockDB := &MockDB{}
	router := SetupGinRouter(mockDB)

	testsCases := []struct {
		Name          string
		RequestMethod string
		Endpoint      string
		Body          io.Reader
		DBErrorMock   error
		ResponseCode  int
	}{
		{
			"HTTP GET request with ACTIVE DB",
			http.MethodGet,
			"/healthz",
			nil,
			nil,
			http.StatusOK,
		},
		{
			"HTTP GET request with INACTIVE DB",
			http.MethodGet,
			"/healthz",
			nil,
			errors.New("mocking database connection failure"),
			http.StatusServiceUnavailable,
		},
		{
			"HTTP GET request WITH params",
			http.MethodGet,
			"/healthz?key=value",
			nil,
			nil,
			http.StatusBadRequest,
		},
		{
			"HTTP GET request with payload",
			http.MethodGet,
			"/healthz",
			strings.NewReader("This request has body"),
			nil,
			http.StatusBadRequest,
		},
	}

	for _, tc := range testsCases {
		t.Run(tc.Name, func(t *testing.T) {
			//Mocking the DB Error to test
			mockDB.Error = tc.DBErrorMock

			_ = mockDB.InitDatabase(config.DatabaseConfig{})

			request := httptest.NewRequest(tc.RequestMethod, tc.Endpoint, tc.Body)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assertStatus(t, response.Code, tc.ResponseCode)
			assertHeader(t, response.Header().Get("Cache-Control"), "no-store, no-cache, must-revalidate;")
		})
	}
}

func TestHealthzRouteWithUnAllowedRequests(t *testing.T) {
	mockDB := &MockDB{}
	router := SetupGinRouter(mockDB)

	testcases := []struct {
		Name          string
		RequestMethod string
	}{
		{
			"HTTP POST request",
			http.MethodPost,
		},
		{
			"HTTP PUT request",
			http.MethodPut,
		},
		{
			"HTTP DELETE request",
			http.MethodDelete,
		},
		{
			"HTTP PATCH request",
			http.MethodPatch,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			request := httptest.NewRequest(tc.RequestMethod, "/healthz", nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assertStatus(t, response.Code, http.StatusMethodNotAllowed)
			assertHeader(t, response.Header().Get("Cache-Control"), "no-store, no-cache, must-revalidate;")
		})
	}
}

func TestForUnhandledRoute(t *testing.T) {
	mockDB := &MockDB{}
	router := SetupGinRouter(mockDB)

	t.Run("Test for unhandled route", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
		assertHeader(t, response.Header().Get("Cache-Control"), "no-store, no-cache, must-revalidate;")
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("Got %d want %d", got, want)
	}
}

func assertHeader(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("Got %s want %s", got, want)
	}
}
