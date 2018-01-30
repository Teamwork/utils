package jsonutil

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMustMarshal(t *testing.T) {
	cases := []struct {
		in   string
		want []byte
	}{
		{`Hello`, []byte(`"Hello"`)},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := MustMarshal(tc.in)
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}
