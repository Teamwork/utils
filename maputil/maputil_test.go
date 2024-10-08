package maputil

import (
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
