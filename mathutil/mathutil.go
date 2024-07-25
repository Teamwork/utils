// Package mathutil provides functions for working with numbers.
package mathutil // import "github.com/teamwork/utils/v2/mathutil"

import (
	"fmt"
	"math"
)

// Round will round the value to the nearest natural number.
//
// .5 will be rounded up.
func Round(f float64) float64 {
	if f < 0 {
		return math.Ceil(f - 0.5)
	}
	return math.Floor(f + 0.5)
}

// RoundPlus will round the value to the given precision.
//
// e.g. RoundPlus(7.258, 2) will return 7.26
func RoundPlus(f float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return Round(f*shift) / shift
}

// CeilPlus will ceil the value to the given precision.
//
// e.g. CeilPlus(123.233333, 2) will return 123.24
func CeilPlus(f float64, precision int) float64 {
	multiplier := math.Pow10(precision)
	return math.Ceil(f*multiplier) / multiplier
}

// Min gets the lowest of two numbers.
func Min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

// Max gets the highest of two numbers.
func Max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

// Limit a value between a lower and upper limit.
func Limit(v, lower, upper float64) float64 {
	return math.Max(math.Min(v, upper), lower)
}

// DivideCeil divides two integers and rounds up, rather than down (which is
// what happens when you do int64/int64).
func DivideCeil(count int64, pageSize int64) int64 {
	return int64(math.Ceil(float64(count) / float64(pageSize)))
}

// IsSignedZero checks if this number is a signed zero (i.e. -0, instead of +0).
func IsSignedZero(f float64) bool {
	return math.Float64bits(f)^uint64(1<<63) == 0
}

// Byte is a float64 where the String() method prints out a human-redable
// description.
type Byte float64

var units = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
var unitsAsBytes = []string{"B", "KB", "MB", "GB", "TB", "PB"}

// String will return the bytes formatted as "mebibytes" (multiples of 1024)
func (b Byte) String() string {
	return b.HumanReadable(1024, units)
}

// StringAsBytes will return the bytes formatted as "bytes" (multiples of 1000)
func (b Byte) StringAsBytes() string {
	return b.HumanReadable(1000, unitsAsBytes)
}

// HumanReadable will take a measurement multiple as well as a slice of formats
// to convert the byte into a human readable format using the given parameters.
func (b Byte) HumanReadable(measurement Byte, format []string) string {
	i := 0
	for ; i < len(units); i++ {
		if b < measurement {
			return fmt.Sprintf("%.1f%s", b, format[i])
		}
		b /= measurement
	}
	return fmt.Sprintf("%.1f%s", b*measurement, format[i-1])
}
