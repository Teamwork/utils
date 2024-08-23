package sliceutil

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/teamwork/test/diff"
)

func TestIntsToString(t *testing.T) {
	cases := []struct {
		in       []int64
		expected string
	}{
		{
			[]int64{1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8},
			"1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8",
		},
		{
			[]int64{-1, -2, -3, -4, -4, -5, -6, -6, -6, -6, -7, -8, -8, -8},
			"-1, -2, -3, -4, -4, -5, -6, -6, -6, -6, -7, -8, -8, -8",
		},
		{
			[]int64{},
			"",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := JoinInt(tc.in)
			if got != tc.expected {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestUniqInt64(t *testing.T) {
	cases := []struct {
		in       []int64
		expected []int64
	}{
		{
			[]int64{1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8},
			[]int64{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			[]int64{1, 3, 8, 3, 8},
			[]int64{1, 3, 8},
		},
		{
			[]int64{1, 2, 3},
			[]int64{1, 2, 3},
		},
		{
			[]int64{},
			nil,
		},
		{
			nil,
			nil,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := UniqInt64(tc.in)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestUniqueMergeSlices(t *testing.T) {
	var tests = []struct {
		in       [][]int64
		expected []int64
	}{
		{
			generate2dintslice([]int64{1, 2, 3}),
			[]int64{1, 2, 3},
		},
		{
			generate2dintslice([]int64{0, 1, 2, 3, -1, -10}),
			[]int64{0, 1, 2, 3, -1, -10},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := UniqueMergeSlices(tc.in)
			if !int64slicesequal(got, tc.expected) {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestUniqString(t *testing.T) {
	var tests = []struct {
		in       []string
		expected []string
	}{
		{
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
		},
		{
			[]string{"a", "b", "c", "a", "b", "n", "a", "aaa", "n", "x"},
			[]string{"a", "b", "c", "n", "aaa", "x"},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := UniqString(tc.in)
			if !stringslicesequal(got, tc.expected) {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestInStringSlice(t *testing.T) {
	tests := []struct {
		list     []string
		find     string
		expected bool
	}{
		{[]string{"hello"}, "hello", true},
		{[]string{"hello"}, "hell", false},
		{[]string{"hello", "world", "test"}, "world", true},
		{[]string{"hello", "world", "test"}, "", false},
		{[]string{}, "", false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := InStringSlice(tc.list, tc.find)
			if got != tc.expected {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestInIntSlice(t *testing.T) {
	tests := []struct {
		list     []int
		find     int
		expected bool
	}{
		{[]int{42}, 42, true},
		{[]int{42}, 4, false},
		{[]int{42, 666, 14159}, 666, true},
		{[]int{42, 666, 14159}, 0, false},
		{[]int{}, 0, false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := InIntSlice(tc.list, tc.find)
			if got != tc.expected {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestInInt64Slice(t *testing.T) {
	tests := []struct {
		list     []int64
		find     int64
		expected bool
	}{
		{[]int64{42}, 42, true},
		{[]int64{42}, 4, false},
		{[]int64{42, 666, 14159}, 666, true},
		{[]int64{42, 666, 14159}, 0, false},
		{[]int64{}, 0, false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := InInt64Slice(tc.list, tc.find)
			if got != tc.expected {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestFilterString(t *testing.T) {
	cases := []struct {
		fun  func(string) bool
		in   []string
		want []string
	}{
		{
			FilterStringEmpty,
			[]string(nil),
			[]string(nil),
		},
		{
			FilterStringEmpty,
			[]string{},
			[]string(nil),
		},
		{
			FilterStringEmpty,
			[]string{"1"},
			[]string{"1"},
		},
		{
			FilterStringEmpty,
			[]string{"", "1", ""},
			[]string{"1"},
		},
		{
			FilterStringEmpty,
			[]string{"", "1", "", "2", "asd", "", "", "", "zx", "", "a"},
			[]string{"1", "2", "asd", "zx", "a"},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := FilterString(tc.in, tc.fun)
			if !reflect.DeepEqual(tc.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func filterIntEmpty(e int64) bool {
	return e != 0
}
func TestFilterInt(t *testing.T) {
	cases := []struct {
		fun  func(int64) bool
		in   []int64
		want []int64
	}{
		{
			filterIntEmpty,
			[]int64(nil),
			[]int64(nil),
		},
		{
			filterIntEmpty,
			[]int64{},
			[]int64(nil),
		},
		{
			filterIntEmpty,
			[]int64{1},
			[]int64{1},
		},
		{
			filterIntEmpty,
			[]int64{0, 1, 0},
			[]int64{1},
		},
		{
			filterIntEmpty,
			[]int64{0, 1, 0, 2, -1, 0, 0, 0, 42, 666, -666, 0, 0, 0},
			[]int64{1, 2, -1, 42, 666, -666},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := FilterInt(tc.in, tc.fun)
			if !reflect.DeepEqual(tc.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestChooseString(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"a"}, "a"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := ChooseString(tt.in)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestRemoveString(t *testing.T) {
	cases := []struct {
		list []string
		item string
		want []string
	}{
		{
			list: []string{"1", "2", "3", "4", "5"},
			item: "3",
			want: []string{"1", "2", "4", "5"},
		},
		{
			list: []string{"1", "2", "3", "4", "5"},
			item: "1",
			want: []string{"2", "3", "4", "5"},
		},
		{
			list: []string{"1", "2", "3", "4", "5"},
			item: "5",
			want: []string{"1", "2", "3", "4"},
		},
		{
			list: []string{"1", "2", "3", "4", "5"},
			item: "6",
			want: []string{"1", "2", "3", "4", "5"},
		},
		{
			list: []string{"2", "1", "2", "2", "2", "5", "2"},
			item: "2",
			want: []string{"1", "5"},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := RemoveString(tc.list, tc.item)
			if !reflect.DeepEqual(tc.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestStringMap(t *testing.T) {
	cases := []struct {
		in   []string
		want []string
		f    func(string) string
	}{
		{
			in:   []string{"a", "b", "c"},
			want: []string{"", "", ""},
			f:    func(string) string { return "" },
		},
		{
			in:   []string{"a", "b", "c"},
			want: []string{"aa", "bb", "cc"},
			f:    func(c string) string { return c + c },
		},
	}
	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			out := StringMap(tc.in, tc.f)
			if !reflect.DeepEqual(tc.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}
