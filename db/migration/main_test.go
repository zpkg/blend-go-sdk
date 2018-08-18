package migration

import (
	"os"
	"testing"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"
)

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	err := db.OpenDefault(db.NewFromEnv())
	if err != nil {
		logger.FatalExit(err)
	}

	os.Exit(m.Run())
}
