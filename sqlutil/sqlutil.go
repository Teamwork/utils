// Package sqlutil provides some helpers for SQL databases.
package sqlutil // import "github.com/teamwork/utils/sqlutil"

import (
	"database/sql/driver"
	"fmt"
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

// Bool add the capability to the sql driver to work directly with bool types.
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
				if raw[0] == 0x1 {
					*b = true
					return nil

				} else if raw[0] == 0x0 {
					*b = false
					return nil
				}
			}

			text = string(raw)

		} else {
			text = v.(string)
		}

		text = strings.ToLower(text)

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
func (b *Bool) Value() (driver.Value, error) {
	if b == nil {
		return nil, fmt.Errorf("boolean not initialized")
	}

	return bool(*b), nil
}

// MarshalText converts the bool to a human readable representation, that is
// also compatible with the JSON format.
func (b *Bool) MarshalText() ([]byte, error) {
	if b == nil {
		return nil, fmt.Errorf("boolean not initialized")
	}

	if *b {
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

	normalized := strings.ToLower(string(text))
	if normalized == "true" || normalized == "1" || normalized == `"true"` {
		*b = true
	} else if normalized == "false" || normalized == "0" || normalized == `"false"` {
		*b = false
	} else {
		return fmt.Errorf("invalid value '%s'", normalized)
	}

	return nil
}
