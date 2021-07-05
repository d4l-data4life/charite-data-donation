package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/go-chi/chi"

	"github.com/d4l-data4life/charite-data-donation/pkg/logging"
	"github.com/d4l-data4life/charite-data-donation/pkg/models"
)

// define error messages
const (
	InvalidDonation = "invalid donation payload"
	ForwardingError = "error forwarding data donation"
)

//NewDonationHandler initializes a new handler
func NewDonationHandler() *DonationHandler {
	return &DonationHandler{}
}

//DonationHandler is the handler responsible for Account operations
type DonationHandler struct {
}

//Routes returns the routes for the DonationHandler
func (e *DonationHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", e.CreateDonation)
	return router
}

//DonationRequest represents the payload for a data donation
type DonationRequest struct {
	PostalCode string `json:"postalCode"`
	RiskCase   int    `json:"riskCase"`
}

func (e *DonationHandler) CreateDonation(w http.ResponseWriter, r *http.Request) {
	donationRequest := &DonationRequest{}
	err := json.NewDecoder(r.Body).Decode(donationRequest)
	if err != nil {
		logging.LogError("Error decoding data donation request payload", err)
		WriteHTTPErrorCode(w, err, http.StatusBadRequest)
		return
	}

	err = ValidateDonation(donationRequest)
	if err != nil {
		WriteHTTPErrorCode(w, err, http.StatusBadRequest)
		return
	}

	donationData, err := json.Marshal(donationRequest)
	if err != nil {
		logging.LogError("Error encoding data donation to json", err)
		WriteHTTPErrorCode(w, err, http.StatusInternalServerError)
		return
	}

	donation := models.Donation{
		Data: donationData,
	}

	err = donation.Create()
	if err != nil {
		WriteHTTPErrorCode(w, err, http.StatusInternalServerError)
		return
	}
	logging.LogInfo("A donation has been received")
}

const PostalCodeRegExp = "^(0[1-9]\\d{3}|[1-9]\\d{4})$"

func ValidateDonation(donation *DonationRequest) error {
	re := regexp.MustCompile(PostalCodeRegExp)
	if !re.MatchString(donation.PostalCode) {
		return errors.New(InvalidDonation)
	}

	if donation.RiskCase < 1 || donation.RiskCase > 5 {
		return errors.New(InvalidDonation)
	}

	return nil
}
