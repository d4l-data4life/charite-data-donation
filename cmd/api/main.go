package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"

	"github.com/d4l-data4life/charite-data-donation/pkg/config"
	"github.com/d4l-data4life/charite-data-donation/pkg/logging"
	"github.com/d4l-data4life/charite-data-donation/pkg/models"
	"github.com/d4l-data4life/charite-data-donation/pkg/server"
)

func main() {
	config.SetupEnv()
	port := viper.GetString("PORT")
	corsHosts := viper.GetString("CORS_HOSTS")

	quitServCh := make(chan struct{})
	quitDBCh := make(chan struct{})
	errorsCh := make(chan config.ErrorMessage)
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT)

	go models.InitializeDB(quitDBCh, errorsCh)

	corsOptions := config.CorsConfig(strings.Split(corsHosts, " "))
	server := server.NewServer(cors.New(corsOptions))
	server.SetupRoutes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(server.Mux(), walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}

	go server.ListenAndServe(quitServCh, errorsCh, port)

	terminateFunc := func(quitDBCh chan struct{}, quitServCh chan struct{}) {
		quitServCh <- struct{}{}
		close(quitServCh)
		quitDBCh <- struct{}{}
		close(quitDBCh)
		time.Sleep(time.Second)
	}

	select {
	case event := <-errorsCh:
		logging.LogInfo("Received error message. Waiting for server to stop...")
		logging.LogError(event.Message, event.Err)
		terminateFunc(quitDBCh, quitServCh)
	case <-termSignal:
		logging.LogInfo("Received terminate signal. Waiting for server to stop...")
		terminateFunc(quitDBCh, quitServCh)
	}
}
