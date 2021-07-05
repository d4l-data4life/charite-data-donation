package config

import (
	"fmt"
	"runtime"

	"github.com/go-chi/cors"
	"github.com/spf13/viper"
)

// Build information. Populated at build-time.
var (
	Name      string = "charite-data-donation"
	Version   string
	Branch    string
	Commit    string
	BuildUser string
	GoVersion = runtime.Version()
)

const (
	//EnvPrefix is a prefix to all ENV variables used in this app
	EnvPrefix = "CHARITE_DATA_DONATION"
	//APIPrefixV1 URL prefix in API version 1
	APIPrefixV1 = "/api/v1"

	//Debug is a flag used to display debug messages
	Debug = false

	// ##### GENERAL VARIABLES

	// DefaultHost default host for the services
	DefaultHost = "localhost"
	// DefaultPort default port the service is served on
	DefaultPort = "4444"
	// DefaultCorsHosts default cors horst for local development
	DefaultCorsHosts = "http://localhost:3333"
	// DefaultDonationURL defines the target to forward donations for the charite
	DefaultForwardURL = ""

	// ##### DATABASE VARIABLES

	// DefaultDBHost default host for the database connection
	DefaultDBHost = "localhost"
	// DefaultDBPort default port for the database connnection
	DefaultDBPort = "5444"
	// DefaultDBName default port for the database connnection
	DefaultDBName = "charite-data-donation"
	// DefaultDBUser default port for the database connnection
	DefaultDBUser = "postgres"
	// DefaultDBPassword default port for the database connnection
	DefaultDBPassword = "postgres"
	// DefaultDBSSLMode default port for the database connnection
	DefaultDBSSLMode = "disable"
)

// ErrorMessage defines the type for the errors channel
type ErrorMessage struct {
	Message string
	Err     error
}

func bindEnvVariable(name string, fallback interface{}) {
	if fallback != "" {
		viper.SetDefault(name, fallback)
	}
	err := viper.BindEnv(name)
	if err != nil {
		//cannot use logging.LogError due to import cycle
		fmt.Printf("Error binding Env Variable: %v", err)
	}
}

//SetupEnv configures app to read ENV variables
func SetupEnv() {
	viper.SetEnvPrefix(EnvPrefix)
	// General
	bindEnvVariable("DEBUG", Debug)
	bindEnvVariable("HOST", DefaultHost)
	bindEnvVariable("PORT", DefaultPort)
	bindEnvVariable("CORS_HOSTS", DefaultCorsHosts)
	bindEnvVariable("FORWARD_URL", DefaultForwardURL)
	// Database
	bindEnvVariable("DB_HOST", DefaultDBHost)
	bindEnvVariable("DB_PORT", DefaultDBPort)
	bindEnvVariable("DB_NAME", DefaultDBName)
	bindEnvVariable("DB_USER", DefaultDBUser)
	bindEnvVariable("DB_PASS", DefaultDBPassword)
	bindEnvVariable("DB_SSL_MODE", DefaultDBSSLMode)
}

// CorsConfig stores default configuration for CORS middleware
func CorsConfig(corsHosts []string) cors.Options {
	return cors.Options{
		AllowedOrigins:   corsHosts,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "X-CSRF-Token"},
		AllowCredentials: true, // header "Access-Control-Allow-Credentials" is not present if this is set to false
		MaxAge:           300,  // Maximum value not ignored by any of major browsers,
		Debug:            viper.GetBool("DEBUG"),
	}
}
