// Package sqlutil provides some helpers for SQL databases.
package sqlutil // import "github.com/teamwork/utils/sqlutil"

import (
	"database/sql/driver"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/teamwork/utils/sliceutil"
)

// IntList expands comma-separated values from a column to []int64, and stores
// []int64 as a comma-separated string.
//
// This is safe for NULL values, in which case it will scan in to IntList(nil).
type IntList []int64

// Value implements the SQL Value function to determine what to store in the DB.
func (l IntList) Value() (driver.Value, error) {
	return sliceutil.JoinInt(l), nil
}

// Scan converts the data returned from the DB into the struct.
func (l *IntList) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	ints := []int64{}
	for _, i := range strings.Split(fmt.Sprintf("%s", v), ",") {
		i = strings.TrimSpace(i)
		if i == "" {
			continue
		}
		in, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			return err
		}
		ints = append(ints, in)
	}
	*l = ints
	return nil
}

// MarshalText converts the data to a human readable representation.
func (l IntList) MarshalText() ([]byte, error) {
	v, err := l.Value()
	return []byte(fmt.Sprintf("%s", v)), err
}

// UnmarshalText parses text in to the Go data structure.
func (l *IntList) UnmarshalText(v []byte) error {
	return l.Scan(v)
}

// FloatList expands comma-separated values from a column to []float64, and
// stores []float64 as a comma-separated string.
//
// This is safe for NULL values, in which case it will scan in to FloatList(nil).
type FloatList []float64

// Value implements the SQL Value function to determine what to store in the DB.
func (l FloatList) Value() (driver.Value, error) {
	return sliceutil.JoinFloat(l), nil
}

// Scan converts the data returned from the DB into the struct.
func (l *FloatList) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	nums := []float64{}
	for _, i := range strings.Split(fmt.Sprintf("%s", v), ",") {
		i = strings.TrimSpace(i)
		if i == "" {
			continue
		}
		in, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return err
		}
		nums = append(nums, in)
	}
	*l = nums
	return nil
}

// MarshalText converts the data to a human readable representation.
func (l FloatList) MarshalText() ([]byte, error) {
	v, err := l.Value()
	return []byte(fmt.Sprintf("%s", v)), err
}

// UnmarshalText parses text in to the Go data structure.
func (l *FloatList) UnmarshalText(v []byte) error {
	return l.Scan(v)
}

// StringList expands comma-separated values from a column to []string, and
// stores []string as a comma-separated string.
//
// Note that this only works for simple strings (e.g. enums), we DO NOT escape
// commas in strings and you will run in to problems.
//
// This is safe for NULL values, in which case it will scan in to
// StringList(nil).
type StringList []string

// Value implements the SQL Value function to determine what to store in the DB.
func (l StringList) Value() (driver.Value, error) {
	return strings.Join(sliceutil.FilterString(l, sliceutil.FilterStringEmpty), ","), nil
}

// Scan converts the data returned from the DB into the struct.
func (l *StringList) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	strs := []string{}
	for _, s := range strings.Split(fmt.Sprintf("%s", v), ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		strs = append(strs, s)
	}
	*l = strs
	return nil
}

// MarshalText converts the data to a human readable representation.
func (l StringList) MarshalText() ([]byte, error) {
	v, err := l.Value()
	return []byte(fmt.Sprintf("%s", v)), err
}

// UnmarshalText parses text in to the Go data structure.
func (l *StringList) UnmarshalText(v []byte) error {
	return l.Scan(v)
}

// Bool converts various column types to a boolean.
//
// Supported types:
//
//   bool
//   int* and float*     0 or 1
//   []byte and string   "1", "true", "0", "false"
//   nil                 defaults to false
type Bool bool

// Scan converts the different types of representation of a boolean in the
// database into a bool type.
func (b *Bool) Scan(src interface{}) error {
	if b == nil {
		return fmt.Errorf("boolean not initialized")
	}

	switch v := src.(type) {
	default:
		return fmt.Errorf("unsupported type %T", src)
	case nil:
		*b = false
	case bool:
		*b = Bool(v)
	case int:
		*b = v != 0
	case int8:
		*b = v != 0
	case int16:
		*b = v != 0
	case int32:
		*b = v != 0
	case int64:
		*b = v != 0
	case uint:
		*b = v != 0
	case uint8:
		*b = v != 0
	case uint16:
		*b = v != 0
	case uint32:
		*b = v != 0
	case uint64:
		*b = v != 0
	case float32:
		*b = v != 0
	case float64:
		*b = v != 0

	case []byte, string:
		var text string
		if raw, ok := v.([]byte); !ok {
			text = v.(string)
		} else if len(raw) == 1 {
			// Handle the bit(1) column type.
			if raw[0] == 1 {
				*b = true
				return nil
			} else if raw[0] == 0 {
				*b = false
				return nil
			}
		} else {
			text = string(raw)
		}

		text = strings.TrimSpace(strings.ToLower(text))
		switch text {
		case "true", "1":
			*b = true
		case "false", "0":
			*b = false
		default:
			return fmt.Errorf("invalid value %q", text)
		}
	}

	return nil
}

// Value converts a bool type into a number to persist it in the database.
func (b Bool) Value() (driver.Value, error) {
	return bool(b), nil
}

// MarshalText converts the data to a JSON-compatible human readable
// representation.
func (b Bool) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%t", b)), nil
}

// UnmarshalText parses text in to the Go data structure.
func (b *Bool) UnmarshalText(text []byte) error {
	if b == nil {
		return fmt.Errorf("boolean not initialized")
	}

	switch strings.TrimSpace(strings.ToLower(string(text))) {
	case "true", "1", `"true"`:
		*b = true
	case "false", "0", `"false"`:
		*b = false
	default:
		return fmt.Errorf("invalid value %q", text)
	}

	return nil
}

// HTML is a string which indicates that the string has been HTML-escaped.
type HTML template.HTML

// Value implements the SQL Value function to determine what to store in the DB.
func (h HTML) Value() (driver.Value, error) {
	return string(h), nil
}

// Scan converts the data returned from the DB into the struct.
func (h *HTML) Scan(v interface{}) error {
	*h = HTML(v.([]byte))
	return nil
}
