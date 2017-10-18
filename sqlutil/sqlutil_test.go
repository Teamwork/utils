package sqlutil

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/teamwork/test"
)

func TestIntListValue(t *testing.T) {
	cases := []struct {
		in   IntList
		want string
	}{
		{*new(IntList), ""},
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
		{"", *new(IntList), ""},
		{"1", IntList{1}, ""},
		{"4, 5", IntList{4, 5}, ""},
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
