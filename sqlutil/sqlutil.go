// Package sqlutil provides some helpers for SQL databases.
package sqlutil // import "github.com/teamwork/utils/v2/sqlutil"

import (
	"database/sql/driver"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/teamwork/utils/v2/sliceutil"
)

// IntList expands comma-separated values from a column to []int64, and stores
// []int64 as a comma-separated string.
//
// This is safe for NULL values, in which case it will scan in to IntList(nil).
type IntList []int64

// Value implements the SQL Value function to determine what to store in the DB.
func (l IntList) Value() (driver.Value, error) {
	return sliceutil.Join(l), nil
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
	return strings.Join(sliceutil.Filter(l, sliceutil.FilterEmpty[string]), ","), nil
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

// Bool add the capability to handle more column types than the usual sql
// driver. The following types are supported when reading the data:
//
//   - int64 and float64 - 0 for false, true otherwise
//   - bool
//   - []byte and string - "1" or "true" for true, and "0" or "false" for false. Also handles the 1 bit cases.
//   - nil - defaults to false
//
// It is also prepared to be encoded and decoded to a human readable format.
type Bool bool

// Scan converts the different types of representation of a boolean in the
// database into a bool type.
func (b *Bool) Scan(src interface{}) error {
	if b == nil {
		return fmt.Errorf("boolean not initialized")
	}

	switch v := src.(type) {
	case int64:
		*b = v != 0
	case float64:
		*b = v != 0
	case bool:
		*b = Bool(v)

	case []byte, string:
		var text string

		if raw, ok := v.([]byte); ok {
			// handle the bit(1) column type
			if len(raw) == 1 {
				switch raw[0] {
				case 0x1:
					*b = true
					return nil

				case 0x0:
					*b = false
					return nil
				}
			}

			text = string(raw)

		} else {
			text = v.(string)
		}

		text = strings.TrimSpace(strings.ToLower(text))

		switch text {
		case "true", "1":
			*b = true
		case "false", "0":
			*b = false
		default:
			return fmt.Errorf("invalid value '%s'", text)
		}

	case nil:
		// nil will be considered false
		*b = false

	default:
		return fmt.Errorf("unsupported format %T", src)
	}

	return nil
}

// Value converts a bool type into a number to persist it in the database.
func (b Bool) Value() (driver.Value, error) {
	return bool(b), nil
}

// MarshalText converts the bool to a human readable representation, that is
// also compatible with the JSON format.
func (b Bool) MarshalText() ([]byte, error) {
	if b {
		return []byte("true"), nil
	}

	return []byte("false"), nil
}

// UnmarshalText parse different types of human representation of the boolean
// and convert it to the bool type. It is also compatible with the JSON format.
func (b *Bool) UnmarshalText(text []byte) error {
	if b == nil {
		return fmt.Errorf("boolean not initialized")
	}

	normalized := strings.TrimSpace(strings.ToLower(string(text)))
	switch normalized {
	case "true", "1", `"true"`:
		*b = true
	case "false", "0", `"false"`:
		*b = false
	default:
		return fmt.Errorf("invalid value '%s'", normalized)
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
