package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/d4l-data4life/charite-data-donation/internal/testutils"
	"github.com/d4l-data4life/charite-data-donation/pkg/handlers"
	"github.com/d4l-data4life/charite-data-donation/pkg/models"
)

func TestRoutesDonations(t *testing.T) {
	router := handlers.NewDonationHandler().Routes()
	assert.NotNil(t, router, "should return a valid router")
	assert.Equal(t, 1, len(router.Routes()), "There should be exactly one routes for this handler")
}

func TestDonationHandler_CreateDonation(t *testing.T) {
	models.InitializeTestDB()
	tests := []struct {
		name       string
		body       handlers.DonationRequest
		statusCode int
	}{
		{"works", handlers.DonationRequest{PostalCode: "12345", RiskCase: 2}, http.StatusOK},
		{"invalid postal code", handlers.DonationRequest{PostalCode: "abcde", RiskCase: 2}, http.StatusBadRequest},
		{"too long postal code", handlers.DonationRequest{PostalCode: "123456", RiskCase: 2}, http.StatusBadRequest},
		{"too low risk case", handlers.DonationRequest{PostalCode: "12345", RiskCase: -1}, http.StatusBadRequest},
		{"too high risk case", handlers.DonationRequest{PostalCode: "12345", RiskCase: 6}, http.StatusBadRequest},
	}
	defer models.GetDB().Close()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &handlers.DonationHandler{}
			request, _ := http.NewRequest("method", "url", testutils.GetRequestPayload(tt.body))
			writer := httptest.NewRecorder()
			e.CreateDonation(writer, request)
			assert.Equal(t, tt.statusCode, writer.Code)
		})
	}
}
