package sqlutil

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/teamwork/test"
)

func TestIntListValue(t *testing.T) {
	cases := []struct {
		in   IntList
		want string
	}{
		{IntList{}, ""},
		{IntList{}, ""},
		{IntList{4, 5}, "4, 5"},
		{IntList{1, 1}, "1, 1"},
		{IntList{1}, "1"},
		{IntList{1, 0, 2}, "1, 0, 2"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out, err := tc.in.Value()
			if err != nil {
				t.Fatal(err)
			}
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestIntListScan(t *testing.T) {
	cases := []struct {
		in      string
		want    IntList
		wantErr string
	}{
		{"", IntList{}, ""},
		{"1", IntList{1}, ""},
		{"4, 5", IntList{4, 5}, ""},
		{"4,   5", IntList{4, 5}, ""},
		{"1, 1", IntList{1, 1}, ""},
		{"1, 0, 2", IntList{1, 0, 2}, ""},
		{"1,0,2", IntList{1, 0, 2}, ""},
		{"1,    0,    2    ", IntList{1, 0, 2}, ""},
		{"1,", IntList{1}, ""},
		{"1,,,,", IntList{1}, ""},
		{",,1,,", IntList{1}, ""},
		{"1,NaN", IntList{}, "strconv.ParseInt"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out := IntList{}
			err := out.Scan(tc.in)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestStringListValue(t *testing.T) {
	cases := []struct {
		in   StringList
		want string
	}{
		{StringList{}, ""},
		{StringList{}, ""},
		{StringList{"4", "5"}, "4,5"},
		{StringList{"1", "1"}, "1,1"},
		{StringList{"€"}, "€"},
		{StringList{"1", "", "1"}, "1,1"},
		{StringList{"لوحة المفاتيح العربية", "xx"}, "لوحة المفاتيح العربية,xx"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out, err := tc.in.Value()
			if err != nil {
				t.Fatal(err)
			}
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestStringListScan(t *testing.T) {
	cases := []struct {
		in      string
		want    StringList
		wantErr string
	}{
		{"", StringList{}, ""},
		{"1", StringList{"1"}, ""},
		{"4, 5", StringList{"4", "5"}, ""},
		{"1, 1", StringList{"1", "1"}, ""},
		{"1,", StringList{"1"}, ""},
		{"1,,,,", StringList{"1"}, ""},
		{",,1,,", StringList{"1"}, ""},
		{"€", StringList{"€"}, ""},
		{"لوحة المفاتيح العربية, xx", StringList{"لوحة المفاتيح العربية", "xx"}, ""},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out := StringList{}
			err := out.Scan(tc.in)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestBoolScan(t *testing.T) {
	cases := []struct {
		in      interface{}
		want    Bool
		wantErr string
	}{
		{[]byte("true"), true, ""},
		{float64(1.0), true, ""},
		{[]byte{0x1}, true, ""},
		{int64(1), true, ""},
		{"true", true, ""},
		{true, true, ""},
		{"1", true, ""},
		{[]byte("false"), false, ""},
		{float64(0.0), false, ""},
		{[]byte{0x0}, false, ""},
		{int64(0), false, ""},
		{"false", false, ""},
		{false, false, ""},
		{"0", false, ""},
		{nil, false, ""},
		{"not a valid bool", false, "invalid value 'not a valid bool'"},
		{time.Time{}, false, "unsupported format time.Time"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			var out Bool
			err := out.Scan(tc.in)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestBoolScanValue(t *testing.T) {
	cases := []struct {
		in   Bool
		want driver.Value
	}{
		{false, false},
		{true, true},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out, err := tc.in.Value()
			if err != nil {
				t.Fatal(err)
			}
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestBoolMarshalText(t *testing.T) {
	cases := []struct {
		in      Bool
		want    []byte
		wantErr string
	}{
		{false, []byte("false"), ""},
		{true, []byte("true"), ""},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out, err := tc.in.MarshalText()
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestBoolUnmarshalText(t *testing.T) {
	cases := []struct {
		in      []byte
		want    Bool
		wantErr string
	}{
		{[]byte("  true  "), true, ""},
		{[]byte(` "true"`), true, ""},
		{[]byte(`  1 `), true, ""},
		{[]byte("false  "), false, ""},
		{[]byte(`"false" `), false, ""},
		{[]byte(` 0 `), false, ""},
		{[]byte(`not a valid bool`), false, "invalid value 'not a valid bool'"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			var out Bool
			err := out.UnmarshalText(tc.in)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}
