package logger

import "testing"

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
