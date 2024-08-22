package jsonutil

import (
	"fmt"
	"reflect"
	"testing"
	"time"
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
				t.Errorf("\nout:  %s\nwant: %s\n", out, tc.want)
			}
		})
	}
}

func TestMustMarshalIndent(t *testing.T) {
	cases := []struct {
		in   map[string]string
		want []byte
	}{
		{map[string]string{"hello": "world", "a": "b"}, []byte("{\n  \"a\": \"b\",\n  \"hello\": \"world\"\n}")},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := MustMarshalIndent(tc.in, "", "  ")
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %s\nwant: %s\n", out, tc.want)
			}
		})
	}
}

func TestMustFormat(t *testing.T) {
	cases := []struct {
		in   []byte
		want []byte
	}{
		{[]byte(`{"hello": "world", "a": "b"}`), []byte("{\n  \"a\": \"b\",\n  \"hello\": \"world\"\n}")},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := MustIndent(tc.in, &map[string]string{}, "", "  ")
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %s\nwant: %s\n", out, tc.want)
			}
		})
	}
}

func TestMustUnmarshal(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var out struct {
			Hello string `json:"hello"`
		}
		MustUnmarshal([]byte(`{"hello":"world"}`), &out)
		if out.Hello != "world" {
			t.Errorf("%#v", out)
		}
	})

	t.Run("panic", func(t *testing.T) {
		defer func() {
			rec := recover()
			if rec == nil {
				t.Error("no panic?")
			}
		}()

		var out struct {
			Hello time.Time `json:"hello"`
		}
		MustUnmarshal([]byte(`{"hello":"world"}`), &out)
	})

}
