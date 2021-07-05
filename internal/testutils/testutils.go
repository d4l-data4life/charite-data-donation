package testutils

import (
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"

	"github.com/go-chi/cors"

	"github.com/d4l-data4life/charite-data-donation/pkg/config"
	"github.com/d4l-data4life/charite-data-donation/pkg/models"
	"github.com/d4l-data4life/charite-data-donation/pkg/server"
)

// GetRequestPayload converts a given object into a reader of that obect as json payload
func GetRequestPayload(payload interface{}) io.Reader {
	bytes, _ := json.Marshal(payload)
	return strings.NewReader(string(bytes))
}

// GetTestMockServer creates the mocked server for tests
func GetTestMockServer() *server.Server {
	config.SetupEnv()
	models.InitializeTestDB()
	corsOptions := config.CorsConfig([]string{"localhost"})
	server := server.NewServer(cors.New(corsOptions))
	server.SetupRoutes()
	return server
}

// RunningTime starts measuring runtime - usage defer Track(RunningTime("label"))
func RunningTime(s string) (string, time.Time) {
	log.Println("Start:	", s)
	return s, time.Now()
}

// Track finishes measuring runtime and prints result - usage defer Track(RunningTime("label"))
func Track(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End:	", s, "took", endTime.Sub(startTime))
}
