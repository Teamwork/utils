// Package sliceutil provides a set if functions for working with slices.
package sliceutil // import "github.com/teamwork/utils/sliceutil"

import (
	"sort"
	"strconv"
	"strings"
)

// JoinInt converts a slice of ints to a comma separated string. Useful for
// inserting into a query without the option of parameterization.
func JoinInt(ints []int64) string {
	var intStr []string
	for _, e := range ints {
		intStr = append(intStr, strconv.Itoa(int(e)))
	}

	return strings.Join(intStr, ", ")
}

// UniqInt64 removes duplicate entries from list. The list does not have to be
// sorted.
func UniqInt64(list []int64) []int64 {
	var unique []int64
	seen := make(map[int64]struct{})
	for _, l := range list {
		if _, ok := seen[l]; !ok {
			seen[l] = struct{}{}
			unique = append(unique, l)
		}
	}
	return unique
}

// UniqString removes duplicate entries from list.
func UniqString(list []string) []string {
	sort.Strings(list)
	var last string
	l := list[:0]
	for _, str := range list {
		if str != last {
			l = append(l, str)
		}
		last = str
	}
	return l
}

// UniqueMergeSlices takes a slice of slices of int64s and returns an unsorted
// slice of unique int64s.
func UniqueMergeSlices(s [][]int64) (result []int64) {
	var m = make(map[int64]bool)

	for _, el := range s {
		for _, i := range el {
			m[i] = true
		}
	}

	for k := range m {
		result = append(result, k)
	}

	return result
}

// CSVtoInt64Slice converts a string of integers to a slice of int64.
func CSVtoInt64Slice(csv string) ([]int64, error) {
	csv = strings.TrimSpace(csv)
	if len(csv) == 0 {
		return []int64(nil), nil
	}

	items := strings.Split(csv, ",")
	ints := make([]int64, len(items))
	for i, item := range items {
		val, err := strconv.Atoi(strings.TrimSpace(item))
		if err != nil {
			return nil, err
		}
		ints[i] = int64(val)
	}

	return ints, nil
}

// InStringSlice reports whether str is within list
func InStringSlice(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}

// InIntSlice reports whether i is within list
func InIntSlice(list []int, i int) bool {
	for _, item := range list {
		if item == i {
			return true
		}
	}
	return false
}

// InInt64Slice reports whether i is within list
func InInt64Slice(list []int64, i int64) bool {
	for _, item := range list {
		if item == i {
			return true
		}
	}
	return false
}

// RepeatString returns a slice with the string s reated n times.
func RepeatString(s string, n int) (r []string) {
	for i := 0; i < n; i++ {
		r = append(r, s)
	}
	return r
}

// Difference returns a new slice with elements that are in "set" but not in
// "others".
func Difference(set []int64, others ...[]int64) []int64 {
	out := []int64{}

	for _, setItem := range set {
		found := false
		for _, o := range others {
			if InInt64Slice(o, setItem) {
				found = true
				break
			}
		}

		if !found {
			out = append(out, setItem)
		}
	}

	return out
}

// Range creates an []int counting at "start" up to (and including) "end".
func Range(start, end int) []int {
	rng := make([]int, end-start+1)
	for i := 0; i < len(rng); i++ {
		rng[i] = start + i
	}
	return rng
}
