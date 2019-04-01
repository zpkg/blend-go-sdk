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
	MaybeInfof(MustNew(), "Infof")
	MaybeDebugf(MustNew(), "Debugf")
	MaybeWarningf(MustNew(), "Warningf")
	MaybeWarning(MustNew(), fmt.Errorf("Warning"))
	MaybeErrorf(MustNew(), "Errorf")
	MaybeError(MustNew(), fmt.Errorf("Error"))
	MaybeFatalf(MustNew(), "Fatalf")
	MaybeFatal(MustNew(), fmt.Errorf("Fatal"))
}
