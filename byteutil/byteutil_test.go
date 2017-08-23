package byteutil

import (
	"fmt"
	"testing"

	"github.com/teamwork/test/diff"
)

func TestToUTF8(t *testing.T) {
	tests := []struct {
		in       []byte
		expected string
	}{
		{[]byte{0x61}, "a"},
		{[]byte{0xd7}, "×"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := ToUTF8(tc.in, "iso-8859-1")
			if got != tc.expected {
				t.Errorf(diff.Cmp(tc.expected, got))
			}
		})
	}
}
