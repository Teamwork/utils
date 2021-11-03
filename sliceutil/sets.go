package sliceutil

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

// Complement returns the complement of the two lists; that is, the first return
// value will contain elements that are only in "a" (and not in "b"), and the
// second return value will contain elements that are only in "b" (and not in
// "a").
func Complement(a, b []int64) (aOnly, bOnly []int64) {
	for _, i := range a {
		if !InInt64Slice(b, i) {
			aOnly = append(aOnly, i)
		}
	}

	for _, i := range b {
		if !InInt64Slice(a, i) {
			bOnly = append(bOnly, i)
		}
	}

	return aOnly, bOnly
}

// Intersection returns the elements common to both a and b
func Intersection(a, b []int64) []int64 {
	inter := []int64{}
	hash := make(map[int64]bool, len(a))
	for _, i := range a {
		hash[i] = false
	}
	for _, i := range b {
		if done, exists := hash[i]; exists && !done {
			inter = append(inter, i)
			hash[i] = true
		}
	}
	return inter
}
