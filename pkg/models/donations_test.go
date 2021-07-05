package models_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/d4l-data4life/charite-data-donation/pkg/models"
)

func TestDonation_Create(t *testing.T) {
	models.InitializeTestDB()
	tests := []struct {
		name    string
		data    interface{}
		wantErr bool
		err     string
	}{
		{"Success", "randomString", false, ""},
	}

	defer models.GetDB().Close()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.data)
			if err != nil {
				t.Errorf("unexpected Error: %v", err)
			}
			Donation := &models.Donation{
				Data: jsonData,
			}
			err = Donation.Create()
			if (err != nil) != tt.wantErr {
				t.Errorf("Donation.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Equal(t, errors.New(tt.err), err, "Test = %v, Msg = Error message should match", tt.name)
			}
		})
	}
}
