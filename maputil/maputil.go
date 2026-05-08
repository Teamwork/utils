// Package maputil provides a set if functions for working with maps.
package maputil // import "github.com/teamwork/utils/v2/maputil"

import (
	"errors"
	"fmt"
)

// Swap the keys and values of a map.
func Swap[T comparable, V comparable](m map[T]V) map[V]T {
	n := make(map[V]T)
	for k, v := range m {
		n[v] = k
	}

	return n
}

var (
	ErrWrongType = errors.New("wrong type")
	ErrNotFound  = errors.New("key not found")
)

// SetValue set value using key path
func SetValue[T any](m map[string]any, v T, keys ...string) error {
	if len(keys) == 1 {
		m[keys[0]] = v
		return nil
	}

	next := m[keys[0]]
	var nextMap map[string]any
	if next == nil {
		nextMap = map[string]any{}
		m[keys[0]] = nextMap
	} else {
		var ok bool
		nextMap, ok = next.(map[string]any)
		if !ok {
			return ErrWrongType
		}
	}

	return SetValue(nextMap, v, keys[1:]...)
}

// GetValue returns the value defined as a key path
func GetValue[T any](m map[string]any, keys ...string) (T, error) {
	var out T

	v, err := getValue(m, keys)
	if err != nil {
		return out, err
	}

	if v == nil {
		return out, ErrNotFound
	}

	vv, ok := v.(T)
	if !ok {
		return out, ErrWrongType
	}

	return vv, nil
}

func getValue(m map[string]any, keys []string) (any, error) {
	if len(keys) == 1 {
		return m[keys[0]], nil
	}

	a := m[keys[0]]

	if m, ok := a.(map[string]any); ok {
		return getValue(m, keys[1:])
	}

	return nil, fmt.Errorf("%w key `%s`", ErrNotFound, keys[1])
}
