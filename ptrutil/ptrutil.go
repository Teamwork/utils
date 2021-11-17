// Package ptrutil contains helper functions related to pointers.
package ptrutil

import (
	"time"
)

// Bool retrieves the pointer of primitive boolean value.
func Bool(b bool) *bool {
	return &b
}

// String retrieves the pointer of primitive string value.
func String(s string) *string {
	return &s
}

// Int retrieves the pointer of primitive int value.
func Int(i int) *int {
	return &i
}

// Int64 retrieves the pointer of primitive int64 value.
func Int64(i int64) *int64 {
	return &i
}

// Float64 retrieves the pointer of primitive float64 value.
func Float64(f float64) *float64 {
	return &f
}

// Time retrieves the pointer of primitive time.Time value.
func Time(t time.Time) *time.Time {
	return &t
}

// BoolOr returns the non-nil value or the given default
func BoolOr(v *bool, d bool) bool {
	if v == nil {
		return d
	}
	return *v
}

// StringOr returns the non-nil value or the given default
func StringOr(v *string, d string) string {
	if v == nil {
		return d
	}
	return *v
}

// IntOr returns the non-nil value or the given default
func IntOr(v *int, d int) int {
	if v == nil {
		return d
	}
	return *v
}

// Int64Or returns the non-nil value or the given default
func Int64Or(v *int64, d int64) int64 {
	if v == nil {
		return d
	}
	return *v
}

// Float64Or returns the non-nil value or the given default
func Float64Or(v *float64, d float64) float64 {
	if v == nil {
		return d
	}

	return *v
}

// TimeOr returns the non-nil value or the given default
func TimeOr(v *time.Time, d time.Time) time.Time {
	if v == nil {
		return d
	}
	return *v
}
