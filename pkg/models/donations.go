package models

import (
	"encoding/json"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/d4l-data4life/charite-data-donation/pkg/logging"
)

// define error messages
const (
	DonationCreationFailure = "failed creating donation in database"
	JSONMarshalFailure      = "failed marshalling JSON data"
)

// Donation model
type Donation struct {
	BaseModel
	Data json.RawMessage `sql:"json"`
}

// Create creates Donation object in the Database
func (donation *Donation) Create() error {
	uid := uuid.NewV4()
	payload, err := donation.Data.MarshalJSON()
	if err != nil {
		logging.LogError(JSONMarshalFailure, err)
		return errors.New(JSONMarshalFailure)
	}
	err = GetDB().Exec("INSERT INTO donations VALUES (?,?,?,?)", uid, time.Now(), time.Now(), payload).Error

	if err != nil {
		logging.LogError(DonationCreationFailure, err)
		return errors.New(DonationCreationFailure)
	}

	return nil
}
