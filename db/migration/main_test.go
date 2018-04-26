package migration

import (
	"log"
	"os"
	"testing"

	"github.com/blend/go-sdk/db"
)

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	err := db.OpenDefault(db.NewFromEnv())
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
