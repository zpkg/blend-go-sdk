package migration

import (
	"log"
	"os"
	"testing"

	"github.com/blend/go-sdk/db"
)

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	connection := db.NewFromEnv()

	err := db.OpenDefault(connection)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
