package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"github.com/d4l-data4life/charite-data-donation/pkg/config"
	"github.com/d4l-data4life/charite-data-donation/pkg/logging"

	// Blank import required by gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var db *gorm.DB

// define general error messages
const (
	ConnectionError = "connection error"
)

// BaseModel defines the basic fields for each other model
type BaseModel struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	if base.ID == uuid.FromStringOrNil("") {
		uuid := uuid.NewV4()
		return scope.SetColumn("ID", uuid)
	}
	return nil
}

// InitializeDB connects to the Database and migrates the schema
func InitializeDB(quit chan struct{}, errors chan config.ErrorMessage) {
	// Attempt to establish DB connection with timeout
	err := ConnectDB()
	for i := 0; i < 5 && err != nil; i++ {
		select {
		case <-quit:
			logging.LogInfo("Database connection not opened")
			return
		default:
			logging.LogWarning("Database Connection failed Trying again in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			err = ConnectDB()
		}
	}
	if err != nil {
		msg := config.ErrorMessage{Message: "Database Connection failed", Err: err}
		select {
		case errors <- msg:
			logging.LogError("Database initialization sent an error to the errors channel", err)
		default:
			logging.LogError("Failed to send on errors channel", nil)
		}
	}

	MigrateDB(true)

	<-quit // this blocks until quit chan receives a value
	CloseDB()
	logging.LogInfo("Database connection closed")
}

//ConnectDB reads environment variables for DB configuration and attempts to open the connection
func ConnectDB() error {
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbName := viper.GetString("DB_NAME")
	dbUser := viper.GetString("DB_USER")
	dbPass := viper.GetString("DB_PASS")
	dbSSLMode := viper.GetString("DB_SSL_MODE")

	connectString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPass, dbSSLMode)

	conn, err := gorm.Open("postgres", connectString)
	if err == nil {
		db = conn
	}
	return err
}

// MigrateDB Executes Migrations on the database
func MigrateDB(debug bool) {
	db := GetDB()
	if debug {
		db = db.Debug()
	}
	db.AutoMigrate(&Donation{})
}

//GetDB returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}

// CloseDB closes the DB connecton
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
