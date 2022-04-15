// Package sqlutil provides some helpers for SQL databases.
package sqlutil // import "github.com/teamwork/utils/sqlutil"

import (
	"database/sql/driver"
	"fmt"
	"html/template"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

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

// Bool add the capability to handle more column types than the usual sql
// driver. The following types are supported when reading the data:
//
//     * int64 and float64 - 0 for false, true otherwise
//     * bool
//     * []byte and string - "1" or "true" for true, and "0" or "false" for false. Also handles the 1 bit cases.
//     * nil - defaults to false
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
	if normalized == "true" || normalized == "1" || normalized == `"true"` { // nolint: gocritic
		*b = true
	} else if normalized == "false" || normalized == "0" || normalized == `"false"` {
		*b = false
	} else {
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

// Interpolate replaces placeholders in the SQL query following Gorp named
// parametes. This should be used only for debugging SQL queries, as it's unsafe
// to execute without prepared statements.
//
// https://github.com/go-gorp/gorp#named-bind-parameters
func Interpolate(query string, args ...interface{}) string {
	if len(args) == 0 {
		return query
	}
	mapParams, ok := args[0].(map[string]interface{})
	if !ok {
		return query
	}

	keys := make([]string, 0, len(mapParams))
	for k := range mapParams {
		keys = append(keys, k)
	}

	// descending sort to replace the longest strings first.
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	for _, key := range keys {
		value := mapParams[key]
		switch v := value.(type) {
		case []string:
			query = strings.ReplaceAll(query, ":"+key, fmt.Sprintf("'%v'", strings.Join(v, "', '")))
		case string:
			query = strings.ReplaceAll(query, ":"+key, fmt.Sprintf("'%v'", value))
		case []bool, []int, []int64, []float64:
			noCommaList := strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%v", value), "["), "]")
			commaList := fmt.Sprintf("%v", strings.Join(strings.Split(noCommaList, " "), ", "))
			query = strings.ReplaceAll(query, ":"+key, commaList)
		case []time.Time:
			strs := make([]string, len(v))
			for i, t := range v {
				strs[i] = t.UTC().Format("2006-01-02 15:04:05")
			}
			query = strings.ReplaceAll(query, ":"+key, fmt.Sprintf("'%v'", strings.Join(strs, "', '")))
		case time.Time:
			query = strings.ReplaceAll(query, ":"+key, fmt.Sprintf("'%v'", v.UTC().Format("2006-01-02 15:04:05")))
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			query = strings.ReplaceAll(query, ":"+key, fmt.Sprintf("%v", value))
		default:
			valueType := reflect.TypeOf(value)
			if valueType.ConvertibleTo(reflect.TypeOf("")) {
				query = strings.ReplaceAll(query, ":"+key, fmt.Sprintf("'%v'", value))
			} else {
				query = strings.ReplaceAll(query, ":"+key, fmt.Sprintf("%v", value))
			}
		}
	}
	return query
}
