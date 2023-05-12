package typeutil_test

import (
	"testing"

	"github.com/teamwork/utils/v2/typeutil"
)

func Test_Default(t *testing.T) {
	tests := map[string]struct {
		in  string
		exp string
	}{
		"empty string returns default": {
			in:  "",
			exp: "default value",
		},
		"non-empty string returns value": {
			in:  "hello there",
			exp: "hello there",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if val := typeutil.Default(test.in, "default value"); val != test.exp {
				t.Fatalf("expected '%s', got '%s'", test.exp, val)
			}
		})
	}
}
