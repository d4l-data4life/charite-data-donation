package handlers_test

import (
	"os"
	"testing"

	"github.com/d4l-data4life/charite-data-donation/pkg/config"
)

// Executed before test runs in this package (fails otherwise)
func TestMain(m *testing.M) {
	config.SetupEnv()
	os.Exit(m.Run())
}
