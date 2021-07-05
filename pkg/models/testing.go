package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// Blank import required by gorm
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// InitializeTestDB connects to an inmemory sqlite for testing
func InitializeTestDB() {
	conn, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	MigrateDB(false)
}
