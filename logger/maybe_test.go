package logger

import (
	"fmt"
	"testing"
)

func TestMaybeNilLogger(t *testing.T) {
	MaybeInfof(nil, "")
	MaybeDebugf(nil, "")
	MaybeWarningf(nil, "")
	MaybeWarning(nil, nil)
	MaybeErrorf(nil, "")
	MaybeError(nil, nil)
	MaybeFatalf(nil, "")
	MaybeFatal(nil, nil)
}

func TestMaybeLogger(t *testing.T) {
	MaybeInfof(New(), "Infof")
	MaybeDebugf(New(), "Debugf")
	MaybeWarningf(New(), "Warningf")
	MaybeWarning(New(), fmt.Errorf("Warning"))
	MaybeErrorf(New(), "Errorf")
	MaybeError(New(), fmt.Errorf("Error"))
	MaybeFatalf(New(), "Fatalf")
	MaybeFatal(New(), fmt.Errorf("Fatal"))
}
