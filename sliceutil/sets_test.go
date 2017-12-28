package sliceutil

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDifference(t *testing.T) {
	cases := []struct {
		inSet    []int64
		inOthers [][]int64
		expected []int64
	}{
		{[]int64{}, [][]int64{}, []int64{}},
		{nil, [][]int64{}, []int64{}},
		{[]int64{}, nil, []int64{}},
		{nil, nil, []int64{}},
		{[]int64{1}, [][]int64{{1}}, []int64{}},
		{[]int64{1, 2, 2, 3}, [][]int64{{1, 2, 2, 3}}, []int64{}},
		{[]int64{1, 2, 2, 3}, [][]int64{{1, 2}, {3}}, []int64{}},
		{[]int64{1, 2}, [][]int64{{1}}, []int64{2}},
		{[]int64{1, 2, 3}, [][]int64{{1}}, []int64{2, 3}},
		{[]int64{1, 2, 3}, [][]int64{{}, {1}}, []int64{2, 3}},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Difference(tc.inSet, tc.inOthers...)
			if !reflect.DeepEqual(tc.expected, out) {
				t.Errorf("\nout:      %#v\nexpected: %#v\n", out, tc.expected)
			}
		})
	}
}

func TestComplement(t *testing.T) {
	type ciTest struct {
		name      string
		a         []int64
		b         []int64
		aExpected []int64
		bExpected []int64
	}
	tests := []ciTest{
		{
			name: "EmptyLists",
		},
		{
			name:      "AOnly",
			a:         []int64{1, 2, 3},
			aExpected: []int64{1, 2, 3},
		},
		{
			name:      "BOnly",
			b:         []int64{1, 2, 3},
			bExpected: []int64{1, 2, 3},
		},
		{
			name: "Equal",
			a:    []int64{1, 2, 3},
			b:    []int64{1, 2, 3},
		},
		{
			name:      "Disjoint",
			a:         []int64{1, 2, 3},
			b:         []int64{5, 6, 7},
			aExpected: []int64{1, 2, 3},
			bExpected: []int64{5, 6, 7},
		},
		{
			name:      "Overlap",
			a:         []int64{1, 2, 3, 4},
			b:         []int64{3, 4, 5, 6},
			aExpected: []int64{1, 2},
			bExpected: []int64{5, 6},
		},
		{
			name:      "Overlap with repeated values",
			a:         []int64{6, 4, 5, 3, 6},
			b:         []int64{2, 1, 4, 3, 1},
			aExpected: []int64{6, 5, 6},
			bExpected: []int64{2, 1, 1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aOnly, bOnly := Complement(test.a, test.b)

			if !reflect.DeepEqual(aOnly, test.aExpected) {
				t.Errorf("aOnly wrong\ngot:  %#v\nwant: %#v\n", aOnly, test.aExpected)
			}
			if !reflect.DeepEqual(bOnly, test.bExpected) {
				t.Errorf("bOnly wrong\ngot:  %#v\nwant: %#v\n", bOnly, test.bExpected)
			}
		})
	}
}

func BenchmarkComplement_equal(b *testing.B) {
	listA := []int64{1, 2, 3}
	listB := []int64{1, 2, 3}

	for n := 0; n < b.N; n++ {
		Complement(listA, listB)
	}
}

func BenchmarkComplement_disjoint(b *testing.B) {
	listA := []int64{1, 2, 3}
	listB := []int64{5, 6, 7}

	for n := 0; n < b.N; n++ {
		Complement(listA, listB)
	}
}

func BenchmarkComplement_overlap(b *testing.B) {
	listA := []int64{1, 2, 3, 4}
	listB := []int64{3, 4, 5, 6}

	for n := 0; n < b.N; n++ {
		Complement(listA, listB)
	}
}
