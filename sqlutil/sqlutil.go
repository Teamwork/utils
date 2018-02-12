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
