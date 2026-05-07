package maputil

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/teamwork/test/diff"
)

func TestSwap(t *testing.T) {
	tests := []struct {
		in       map[string]string
		expected map[string]string
	}{
		{map[string]string{"a": "b"}, map[string]string{"b": "a"}},
		{map[string]string{"a": "b", "c": "d"}, map[string]string{"b": "a", "d": "c"}},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := Swap(tc.in)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestGetValue(t *testing.T) {
	m := map[string]any{
		"a": int64(1),
		"b": "2",
		"c": float64(3.1),
		"d": map[string]any{
			"a": int64(4),
			"b": "5",
			"c": float64(6.1),
		},
	}

	outInt64, err := GetValue[int64](m, "a")
	if err != nil {
		t.Fatal(err)
	}
	if outInt64 != int64(1) {
		t.Fatalf("expected 1, got %v", outInt64)
	}

	outStr, err := GetValue[string](m, "b")
	if err != nil {
		t.Fatal(err)
	}
	if outStr != "2" {
		t.Fatalf("expected \"2\", got %v", outStr)
	}

	outFloat, err := GetValue[float64](m, "c")
	if err != nil {
		t.Fatal(err)
	}
	if outFloat != 3.1 {
		t.Fatalf("expected 3.1, got %v", outFloat)
	}

	outInt64, err = GetValue[int64](m, "d", "a")
	if err != nil {
		t.Fatal(err)
	}
	if outInt64 != int64(4) {
		t.Fatalf("expected 4 got %v", outInt64)
	}

	_, err = GetValue[string](m, "a")
	if !errors.Is(err, ErrWrongType) {
		t.Fatalf("error wrong type expected, got %v", err)
	}

	_, err = GetValue[string](m, "d", "a")
	if !errors.Is(err, ErrWrongType) {
		t.Fatalf("error wrong type expected, got %v", err)
	}

	_, err = GetValue[string](m, "e")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("error not found expected, got %v", err)
	}

	_, err = GetValue[string](m, "d", "e")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("error not found expected, got %v", err)
	}

	_, err = GetValue[string](m, "e", "d")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("error not found expected, got %v", err)
	}

	_, err = GetValue[string](m, "d", "c", "a")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("error not found expected, got %v", err)
	}

	_, err = GetValue[string](m, "a", "")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("error not found expected, got %v", err)
	}
}
