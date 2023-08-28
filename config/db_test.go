package config

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestConnect(t *testing.T) {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("error loading env variables %s", err.Error())
	}

	// Run the Connect function
	Connect()

	// Check if the DB variable is not nil
	if DB == nil {
		t.Fatalf("DB connection is nil")
	}

	// Check if the DB connection is valid
	err = DB.Raw("SELECT 1").Error
	if err != nil {
		t.Fatalf("DB connection is invalid: %s", err.Error())
	}
}
