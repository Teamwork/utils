package sliceutil

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/teamwork/test/diff"
)

func TestJoin(t *testing.T) {
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
			got := Join(tc.in)
			if got != tc.expected {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestJoinWith(t *testing.T) {
	cases := []struct {
		in       []int64
		delim    string
		expected string
	}{
		{
			[]int64{1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8},
			" || ",
			"1 || 2 || 3 || 4 || 4 || 5 || 6 || 6 || 6 || 6 || 7 || 8 || 8 || 8",
		},
		{
			[]int64{1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8},
			",",
			"1,2,3,4,4,5,6,6,6,6,7,8,8,8",
		},
		{
			[]int64{-1, -2, -3, -4, -4, -5, -6, -6, -6, -6, -7, -8, -8, -8},
			" || ",
			"-1 || -2 || -3 || -4 || -4 || -5 || -6 || -6 || -6 || -6 || -7 || -8 || -8 || -8",
		},
		{
			[]int64{-1, -2, -3, -4, -4, -5, -6, -6, -6, -6, -7, -8, -8, -8},
			",",
			"-1,-2,-3,-4,-4,-5,-6,-6,-6,-6,-7,-8,-8,-8",
		},
		{
			[]int64{},
			" || ",
			"",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := JoinWith(tc.in, tc.delim)
			if got != tc.expected {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestUniq_Int64(t *testing.T) {
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
			got := Unique(tc.in)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestUniq_String(t *testing.T) {
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
			got := Unique(tc.in)
			if !stringslicesequal(got, tc.expected) {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestMergeUnique_Int64(t *testing.T) {
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
			got := MergeUnique(tc.in)
			if !int64slicesequal(got, tc.expected) {
				t.Log("IN", tc.in)
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func int64slicesequal(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}

	for _, ia := range a {
		var found bool
		for _, ib := range b {
			if ib == ia {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func stringslicesequal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for _, ia := range a {
		var found bool
		for _, ib := range b {
			if ib == ia {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func generate2dintslice(in []int64) [][]int64 {
	var (
		result    [][]int64
		processed = map[int]struct{}{}
		loops     = int(rand.Int63n(int64(len(in)*2)) + 1)
	)

	for len(processed) < len(in) {
		var s []int64
		for i := 0; i < loops; i++ {
			idx := rand.Intn(len(in))
			processed[idx] = struct{}{}
			s = append(s, in[idx])
		}
		result = append(result, s)
	}

	return result
}

func TestCSVtoInt64Slice(t *testing.T) {
	tests := []struct {
		in          string
		expected    []int64
		expectedErr error
	}{
		{
			"1,2,3",
			[]int64{1, 2, 3},
			nil,
		},
		{
			"",
			[]int64(nil),
			nil,
		},
		{
			"1,				2, \n3",
			[]int64{1, 2, 3},
			nil,
		},
		{
			"1,				2,nope",
			[]int64(nil),
			errors.New("invalid syntax"),
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got, err := CSVtoInt64Slice(tc.in)

			if err != nil {
				if numErrorer, ok := err.(*strconv.NumError); ok {
					err = numErrorer.Err
				}
			}

			if err != tc.expectedErr && err.Error() != tc.expectedErr.Error() {
				t.Error(diff.Cmp(tc.expectedErr.Error(), err.Error()))
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestItemInSlice_String(t *testing.T) {
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
			got := Contains(tc.list, tc.find)
			if got != tc.expected {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestInFoldedStringSlice(t *testing.T) {
	tests := []struct {
		list     []string
		find     string
		expected bool
	}{
		{[]string{"hello"}, "hello", true},
		{[]string{"HELLO"}, "hello", true},
		{[]string{"hello"}, "HELLO", true},
		{[]string{"hello"}, "hell", false},
		{[]string{"hello", "world", "test"}, "world", true},
		{[]string{"hello", "world", "test"}, "", false},
		{[]string{}, "", false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := InFoldedStringSlice(tc.list, tc.find)
			if got != tc.expected {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestItemInSlice_Int(t *testing.T) {
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
			got := Contains(tc.list, tc.find)
			if got != tc.expected {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestItemInSlice_Int64(t *testing.T) {
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
			got := Contains(tc.list, tc.find)
			if got != tc.expected {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}

func TestRange(t *testing.T) {
	cases := []struct {
		start, end int
		want       []int
	}{
		{1, 5, []int{1, 2, 3, 4, 5}},
		{0, 5, []int{0, 1, 2, 3, 4, 5}},
		{-2, 5, []int{-2, -1, 0, 1, 2, 3, 4, 5}},
		{-5, -1, []int{-5, -4, -3, -2, -1}},
		{100, 105, []int{100, 101, 102, 103, 104, 105}},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v-%v", tc.start, tc.end), func(t *testing.T) {
			out := Range(tc.start, tc.end)
			if !reflect.DeepEqual(tc.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestFilter_String(t *testing.T) {
	cases := []struct {
		fun  func(string) bool
		in   []string
		want []string
	}{
		{
			FilterEmpty[string],
			[]string(nil),
			[]string(nil),
		},
		{
			FilterEmpty[string],
			[]string{},
			[]string(nil),
		},
		{
			FilterEmpty[string],
			[]string{"1"},
			[]string{"1"},
		},
		{
			FilterEmpty[string],
			[]string{"", "1", ""},
			[]string{"1"},
		},
		{
			FilterEmpty[string],
			[]string{"", "1", "", "2", "asd", "", "", "", "zx", "", "a"},
			[]string{"1", "2", "asd", "zx", "a"},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Filter(tc.in, tc.fun)
			if !reflect.DeepEqual(tc.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestFilter_Int(t *testing.T) {
	cases := []struct {
		fun  func(int64) bool
		in   []int64
		want []int64
	}{
		{
			FilterEmpty[int64],
			[]int64(nil),
			[]int64(nil),
		},
		{
			FilterEmpty[int64],
			[]int64{},
			[]int64(nil),
		},
		{
			FilterEmpty[int64],
			[]int64{1},
			[]int64{1},
		},
		{
			FilterEmpty[int64],
			[]int64{0, 1, 0},
			[]int64{1},
		},
		{
			FilterEmpty[int64],
			[]int64{0, 1, 0, 2, -1, 0, 0, 0, 42, 666, -666, 0, 0, 0},
			[]int64{1, 2, -1, 42, 666, -666},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Filter(tc.in, tc.fun)
			if !reflect.DeepEqual(tc.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestChoose_String(t *testing.T) {
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
			out := Choose(tt.in)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestRemove_String(t *testing.T) {
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
			out := Remove(tc.list, tc.item)
			if !reflect.DeepEqual(tc.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestMap_String(t *testing.T) {
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
			out := Map(tc.in, tc.f)
			if !reflect.DeepEqual(tc.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestInterfaceSliceTo(t *testing.T) {
	{
		src := []interface{}{"1", "2", "3", "4", "5"}
		want := []string{"1", "2", "3", "4", "5"}
		input := []string{""}

		result := InterfaceSliceTo(src, input)
		if !reflect.DeepEqual(result, want) {
			t.Fatalf("want %+v(%T),\tgot %+v(%T)", want, want, result, result)
		}
	}

	{
		src := []interface{}{"1", "2", "3", "4", "5"}
		want := []string{"1", "2", "3", "4", "5"}
		input := []string{}

		result := InterfaceSliceTo(src, input)
		if !reflect.DeepEqual(result, want) {
			t.Fatalf("want %+v(%T),\tgot %+v(%T)", want, want, result, result)
		}
	}

	{
		src := []interface{}{1, 2, 3, 4, 5}
		want := []int{1, 2, 3, 4, 5}
		input := []int{0, 0}

		result := InterfaceSliceTo(src, input)
		if !reflect.DeepEqual(result, want) {
			t.Fatalf("want %+v(%T),\tgot %+v(%T)", want, want, result, result)
		}
	}

	{
		src := []interface{}{1, 2, 3, 4, 5}
		want := []int64{1, 2, 3, 4, 5}
		input := []int64{0, 0}

		result := InterfaceSliceTo(src, input)
		if !reflect.DeepEqual(result, want) {
			t.Fatalf("want %+v(%T),\tgot %+v(%T)", want, want, result, result)
		}
	}
}

func TestToAnySlice(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	out := ToAnySlice(in)
	for i := range out {
		if !reflect.DeepEqual(out[i], in[i]) {
			t.Errorf("want %[1]v(%[1]T),\tgot %[2]v(%[2]T)", in[i], out[i])
		}
	}
}

func TestValues(t *testing.T) {
	type testStruct struct {
		Name string
		Age  int
	}

	cases := []struct {
		in       []testStruct
		expected []string
	}{
		{
			[]testStruct{
				{Name: "a", Age: 1},
				{Name: "b", Age: 2},
				{Name: "c", Age: 3},
			},
			[]string{"a", "b", "c"},
		},
		{
			[]testStruct{
				{Name: "a", Age: 1},
				{Name: "b", Age: 2},
				{Name: "c", Age: 3},
				{Name: "d", Age: 4},
			},
			[]string{"a", "b", "c", "d"},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := Values(tc.in, func(t testStruct) string {
				return t.Name
			})
			if !reflect.DeepEqual(got, tc.expected) {
				t.Error(diff.Cmp(tc.expected, got))
			}
		})
	}
}
