package ptrutil

import (
	"fmt"
	"testing"

	"github.com/teamwork/test/diff"
)

func TestDereference_Int(t *testing.T) {
	tests := []struct {
		ptr      *int
		expected int
	}{{
		ptr:      func() *int { i := 1; return &i }(),
		expected: 1,
	}, {
		ptr:      func() *int { i := 10; return &i }(),
		expected: 10,
	}, {
		expected: 0,
	}}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			if got := Dereference(test.ptr); got != test.expected {
				t.Error(diff.Cmp(test.expected, got))
			}
		})
	}
}

func TestDereference_String(t *testing.T) {
	tests := []struct {
		ptr      *string
		expected string
	}{{
		ptr:      func() *string { i := "hello"; return &i }(),
		expected: "hello",
	}, {
		ptr:      func() *string { i := "world"; return &i }(),
		expected: "world",
	}, {
		expected: "",
	}}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			if got := Dereference(test.ptr); got != test.expected {
				t.Error(diff.Cmp(test.expected, got))
			}
		})
	}
}
