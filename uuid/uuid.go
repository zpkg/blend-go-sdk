package uuid

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strings"

	"github.com/blend/go-sdk/exception"
)

var (
	byteGroups = []int{8, 4, 4, 4, 12}

	byteGroupSeparatorOffsets = []int{8, 12, 16, 20}

	hextable = [16]byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
	}
)

// Empty returns an empty uuid block.
func Empty() UUID {
	return UUID(make([]byte, 16))
}

// UUID represents a unique identifier conforming to the RFC 4122 standard.
// UUIDs are a fixed 128bit (16 byte) binary blob.
type UUID []byte

// Equals returns if a uuid equals another uuid.
func (uuid UUID) Equals(other UUID) bool {
	if uuid == nil || other == nil {
		return false
	}
	if len(uuid) != len(other) {
		return false
	}
	for index := 0; index < len(uuid); index++ {
		if uuid[index] != other[index] {
			return false
		}
	}
	return true
}

// ToFullString returns a "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" hex representation of a uuid.
func (uuid UUID) ToFullString() string {
	if len(uuid) == 0 {
		return ""
	}
	b := []byte(uuid)
	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		b[:4], b[4:6], b[6:8], b[8:10], b[10:],
	)
}

// ToShortString returns a hex representation of the uuid.
func (uuid UUID) ToShortString() string {
	b := []byte(uuid)
	return fmt.Sprintf("%x", b[:])
}

// String is an alias for `ToShortString`.
func (uuid UUID) String() string {
	return uuid.ToShortString()
}

// Version returns the version byte of a uuid.
func (uuid UUID) Version() byte {
	return uuid[6] >> 4
}

// Format allows for conditional expansion in printf statements
// based on the token and flags used.
func (uuid UUID) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, uuid.ToFullString())
			return
		}
		io.WriteString(s, uuid.ToShortString())
	case 's':
		io.WriteString(s, uuid.ToShortString())
	case 'q':
		fmt.Fprintf(s, "%b", uuid.Version())
	}
}

// IsV4 returns true iff uuid has version number 4, variant number 2, and length 16 bytes
func (uuid UUID) IsV4() bool {
	if len(uuid) != 16 {
		return false
	}
	// check that version number is 4
	if (uuid[6]&0xf0)^0x40 != 0 {
		return false
	}
	// check that variant is 2
	return (uuid[8]&0xc0)^0x80 == 0
}

// MarshalJSON marshals a uuid as json.
func (uuid UUID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + uuid.ToFullString() + "\""), nil
}

// UnmarshalJSON unmarshals a uuid from json.
func (uuid *UUID) UnmarshalJSON(corpus []byte) error {
	if len(*uuid) == 0 {
		(*uuid) = Empty()
	}
	raw := strings.TrimSpace(string(corpus))
	raw = strings.TrimPrefix(raw, "\"")
	raw = strings.TrimSuffix(raw, "\"")
	return ParseExisting(uuid, raw)
}

// MarshalYAML marshals a uuid as yaml.
func (uuid UUID) MarshalYAML() (interface{}, error) {
	return "\"" + uuid.ToFullString() + "\"", nil
}

// UnmarshalYAML unmarshals a uuid from yaml.
func (uuid *UUID) UnmarshalYAML(unmarshaler func(interface{}) error) error {
	if len(*uuid) == 0 {
		(*uuid) = Empty()
	}

	var corpus string
	if err := unmarshaler(&corpus); err != nil {
		return err
	}

	raw := strings.TrimSpace(string(corpus))
	raw = strings.TrimPrefix(raw, "\"")
	raw = strings.TrimSuffix(raw, "\"")
	return ParseExisting(uuid, raw)
}

// Scan scans a uuid from a db value.
func (uuid *UUID) Scan(src interface{}) error {
	if len(*uuid) == 0 {
		(*uuid) = Empty()
	}
	switch src.(type) {
	case string:
		return ParseExisting(uuid, src.(string))
	case []byte:
		return ParseExisting(uuid, string(src.([]byte)))
	}
	return exception.New(exception.Class("uuid: invalid scan source")).WithMessagef("scan type: %T", src)
}

// Value returns a sql driver value.
func (uuid UUID) Value() (driver.Value, error) {
	if uuid == nil || len(uuid) == 0 {
		return nil, nil
	}
	return uuid.ToFullString(), nil
}
