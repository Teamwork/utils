package mathutil

import (
	"sort"
	"testing"

	"github.com/teamwork/test/diff"
)

func TestComplementsInt(t *testing.T) {
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
			aExpected: []int64{5, 6},
			bExpected: []int64{1, 2},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aOnly, bOnly := ComplementsInt(test.a, test.b)
			sort.Slice(aOnly, func(i, j int) bool { return aOnly[i] < aOnly[j] })
			sort.Slice(bOnly, func(i, j int) bool { return bOnly[i] < bOnly[j] })
			if d := diff.MarshalJSONDiff(test.aExpected, aOnly); d != "" {
				t.Errorf("A differs:\n%s", d)
			}
			if d := diff.MarshalJSONDiff(test.bExpected, bOnly); d != "" {
				t.Errorf("B differs:\n%s", d)
			}
		})
	}
}

func BenchmarkComplementsInt_equal(b *testing.B) {
	listA := []int64{1, 2, 3}
	listB := []int64{1, 2, 3}

	for n := 0; n < b.N; n++ {
		ComplementsInt(listA, listB)
	}
}

func BenchmarkComplementsInt_disjoint(b *testing.B) {
	listA := []int64{1, 2, 3}
	listB := []int64{5, 6, 7}

	for n := 0; n < b.N; n++ {
		ComplementsInt(listA, listB)
	}
}

func BenchmarkComplementsInt_overlap(b *testing.B) {
	listA := []int64{1, 2, 3, 4}
	listB := []int64{3, 4, 5, 6}

	for n := 0; n < b.N; n++ {
		ComplementsInt(listA, listB)
	}
}
