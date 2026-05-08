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

func TestSetValue(t *testing.T) {
	m := map[string]any{}

	// simple int set
	err := SetValue(m, 10, "a")
	if err != nil {
		t.Fatal(err)
	}
	if m["a"] == nil {
		t.Fatal("expected key 'a'")
	}
	if v, ok := m["a"].(int); !ok || v != 10 {
		t.Fatalf("expected 10, got %v", v)
	}

	// simple string set
	str := "mystring"
	err = SetValue(m, str, "b")
	if err != nil {
		t.Fatal(err)
	}
	if m["b"] == nil {
		t.Fatal("expected key 'b'")
	}
	if v, ok := m["b"].(string); !ok || v != str {
		t.Fatalf("expected '%v', got %v", str, v)
	}

	// nested set
	err = SetValue(m, 2, "c", "d")
	if err != nil {
		t.Fatal(err)
	}
	if m["c"] == nil {
		t.Fatal("expected key 'b'")
	}
	nested, ok := m["c"].(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", nested)
	}
	if nested["d"] == nil {
		t.Fatal("expected key 'd'")
	}
	if v, ok := nested["d"].(int); !ok || v != 2 {
		t.Fatalf("expected 2, got %v", v)
	}

	// replace
	err = SetValue(map[string]any{"a": 10}, 10, "a")
	if err != nil {
		t.Fatal(err)
	}
	if m["a"] == nil {
		t.Fatal("expected key 'a'")
	}
	if v, ok := m["a"].(int); !ok || v != 10 {
		t.Fatalf("expected 10, got %v", v)
	}

	// replace tip
	m = map[string]any{"a": map[string]any{"b": 1}}
	err = SetValue(m, 10, "a", "b")
	if err != nil {
		t.Fatal(err)
	}
	if v, _ := m["a"].(map[string]any)["b"].(int); v != 10 {
		t.Fatalf("expected 10, got %v", v)
	}

	// crash when mid key is not a map
	m = map[string]any{"a": map[string]any{"b": 1}}
	err = SetValue(m, 10, "a", "b", "c")
	if err != ErrWrongType {
		t.Fatal(err)
	}
}
