package secrets

import (
	"github.com/blend/go-sdk/assert"
	"testing"
)

func TestTransitCreateOptConfig(t *testing.T) {
	a := assert.New(t)
	cfg := CreateTransitKeyConfig{
		Derived: true,
		Convergent: true,
	}

	empty := &CreateTransitKeyConfig{}

	TKCreateOptConfig(cfg)(empty)

	a.True(empty.Derived)
	a.True(empty.Convergent)
}

func TestTransitUpdateOptConfig(t *testing.T) {
	a := assert.New(t)
	cfg := UpdateTransitKeyConfig{
		Exportable: true,
		AllowPlaintextBackup: true,
	}

	empty := &UpdateTransitKeyConfig{}

	TKUpdateOptConfig(cfg)(empty)

	a.True(empty.Exportable)
	a.True(empty.AllowPlaintextBackup)
}

func TestTransitCreateOptMisc(t *testing.T) {
	a := assert.New(t)
	empty := &CreateTransitKeyConfig{}

	TKCreateOptDerived()(empty)
	a.True(empty.Derived)

	empty.Derived = false

	TKCreateOptConvergent()(empty)
	a.True(empty.Derived)
	a.True(empty.Convergent)

	TKCreateOptAllowPlaintextBackup()(empty)
	a.True(empty.AllowPlaintextBackup)

	TKCreateOptExportable()(empty)
	a.True(empty.Exportable)

	err := TKCreateOptType("not a real type")(empty)
	a.NotNil(err)

	err = TKCreateOptType(TypeCHACHA20POLY1305)(empty)
	a.Nil(err)
	a.Equal(TypeCHACHA20POLY1305, empty.Type)
}

func TestTransitUpdateOptMisc(t *testing.T) {
	a := assert.New(t)
	empty := &UpdateTransitKeyConfig{}

	TKUpdateOptDeletionAllowed(true)(empty)
	a.True(*empty.DeletionAllowed)

	TKUpdateOptDeletionAllowed(false)(empty)
	a.False(*empty.DeletionAllowed)

	TKUpdateOptAllowPlaintextBackup()(empty)
	a.True(empty.AllowPlaintextBackup)

	TKUpdateOptExportable()(empty)
	a.True(empty.Exportable)

	TKUpdateOptMinDecryptionVersion(4)(empty)
	a.Equal(4, empty.MinDecryptionVersion)

	TKUpdateOptMinEncryptionnVersion(5)(empty)
	a.Equal(5, empty.MinEncryptionVersion)
}
