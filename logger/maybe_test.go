package logger

import (
	"fmt"
	"testing"
)

func TestMaybeNilLogger(t *testing.T) {
	MaybeInfof(nil, "")
	MaybeSyncInfof(nil, "")
	MaybeDebugf(nil, "")
	MaybeSyncDebugf(nil, "")
	MaybeWarningf(nil, "")
	MaybeSyncWarningf(nil, "")
	MaybeWarning(nil, nil)
	MaybeSyncWarning(nil, nil)
	MaybeErrorf(nil, "")
	MaybeSyncErrorf(nil, "")
	MaybeError(nil, nil)
	MaybeSyncError(nil, nil)
	MaybeFatalf(nil, "")
	MaybeSyncFatalf(nil, "")
	MaybeFatal(nil, nil)
	MaybeSyncFatal(nil, nil)
}
func TestMaybeLogger(t *testing.T) {
	MaybeInfof(New(), "Infof")
	MaybeSyncInfof(New(), "SyncInfof")
	MaybeDebugf(New(), "Debugf")
	MaybeSyncDebugf(New(), "SyncDebugf")
	MaybeWarningf(New(), "Warningf")
	MaybeSyncWarningf(New(), "SyncWarningf")
	MaybeWarning(New(), fmt.Errorf("Warning"))
	MaybeSyncWarning(New(), fmt.Errorf("SyncWarning"))
	MaybeErrorf(New(), "Errorf")
	MaybeSyncErrorf(New(), "SyncErrorf")
	MaybeError(New(), fmt.Errorf("Error"))
	MaybeSyncError(New(), fmt.Errorf("SyncError"))
	MaybeFatalf(New(), "Fatalf")
	MaybeSyncFatalf(New(), "SyncFatalf")
	MaybeFatal(New(), fmt.Errorf("Fatal"))
	MaybeSyncFatal(New(), fmt.Errorf("SyncFatal"))
}
