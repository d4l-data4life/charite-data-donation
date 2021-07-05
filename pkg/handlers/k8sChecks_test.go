package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/d4l-data4life/charite-data-donation/pkg/handlers"
	"github.com/d4l-data4life/charite-data-donation/pkg/models"
)

const (
	livenessURL  = "/checks/liveness"
	readinessURL = "/checks/readiness"
)

func TestRoutesCheck(t *testing.T) {
	router := handlers.NewChecksHandler().Routes()
	assert.NotNil(t, router, "should return a valid router")
	assert.Equal(t, 2, len(router.Routes()), "There should be exactly two routes for this handler")
}

func TestCheckLiveness(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, livenessURL, nil)
	response := httptest.NewRecorder()
	handlers.NewChecksHandler().Liveness(response, request)
	assert.Equal(t, 200, response.Code)
}

func TestCheckReadiness(t *testing.T) {
	models.InitializeTestDB()
	request, _ := http.NewRequest(http.MethodGet, readinessURL, nil)
	response := httptest.NewRecorder()
	handlers.NewChecksHandler().Readiness(response, request)
	assert.Equal(t, 200, response.Code)
}

func TestCheckReadinessFailure(t *testing.T) {
	// Open and Close DB connection to simulate broken connection
	models.InitializeTestDB()
	models.GetDB().Close()
	request, _ := http.NewRequest(http.MethodGet, readinessURL, nil)
	response := httptest.NewRecorder()
	handlers.NewChecksHandler().Readiness(response, request)
	assert.Equal(t, 500, response.Code)
}
