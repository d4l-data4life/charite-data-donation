package handlers

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/d4l-data4life/charite-data-donation/pkg/logging"
	"github.com/d4l-data4life/charite-data-donation/pkg/models"
)

//ChecksHandler is the handler responsible for k8s checks
type ChecksHandler struct {
}

//Routes returns the routes for the ChecksHandler
func (e *ChecksHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/liveness", e.Liveness)
	router.Get("/readiness", e.Readiness)
	return router
}

//NewChecksHandler initializes a new handler
func NewChecksHandler() *ChecksHandler {
	return &ChecksHandler{}
}

//Liveness is a check that describes if the application has started
func (e *ChecksHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	WriteHTTPCode(w, http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		logging.LogError("Error writing OK to response body", err)
	}
}

//Readiness is a check if application can handle requests
func (e *ChecksHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	if err := models.GetDB().DB().Ping(); err != nil {
		WriteHTTPErrorCode(w, err, http.StatusInternalServerError)
		return
	}

	WriteHTTPCode(w, http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		logging.LogError("Error writing OK to response body", err)
	}
}
