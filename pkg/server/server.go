package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	"github.com/d4l-data4life/charite-data-donation/pkg/config"
	"github.com/d4l-data4life/charite-data-donation/pkg/handlers"
	"github.com/d4l-data4life/charite-data-donation/pkg/logging"
)

// Server wraps a chi router (chi.Mux)
type Server struct {
	mux *chi.Mux
}

func newMux(cors *cors.Cors) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		cors.Handler, // Set Access-Control-Allow-Origin header
		middleware.RequestID,
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.Recoverer,       // Recover from panics without crashing server
		middleware.StripSlashes,
	)
	return mux
}

// NewServer creates a router with routes setup
func NewServer(cors *cors.Cors) *Server {
	server := Server{mux: newMux(cors)}
	return &server
}

// Mux returns the chi router
func (s *Server) Mux() *chi.Mux {
	return s.mux
}

// SetupRoutes adds all routes that the server should listen to
func (s *Server) SetupRoutes() {
	ch := handlers.NewChecksHandler()
	donationHandler := handlers.NewDonationHandler()
	s.Mux().Mount("/checks", ch.Routes())

	s.Mux().Route(config.APIPrefixV1, func(r chi.Router) {
		r.Mount("/donations", donationHandler.Routes())
	})
}

// ListenAndServe starts the server
func (s *Server) ListenAndServe(quit chan struct{}, errors chan config.ErrorMessage, port string) {
	go func() {
		listenAddress := net.JoinHostPort("", port)
		logging.LogInfo(fmt.Sprintf("Listeninig on %s\n", listenAddress))
		if err := http.ListenAndServe(listenAddress, s.mux); err != nil {
			msg := config.ErrorMessage{Message: fmt.Sprintf("Could not listen on port %s", port), Err: err}
			select {
			case errors <- msg:
				logging.LogInfo("Sent on errors channel")
			default:
				logging.LogInfo("Failed to send on errors channel")
			}
		}
	}()

	<-quit // this blocks until quit chan receives a value
	logging.LogInfo("Server has been shutdown")
}
