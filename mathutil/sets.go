package mathutil

// ComplementsInt takes two int64 slices, and returns the complements of the
// two lists.
func ComplementsInt(a, b []int64) (aOnly, bOnly []int64) {
	aMap := make(map[int64]struct{}, len(a))
	for _, i := range a {
		aMap[i] = struct{}{}
	}
	for _, i := range b {
		if _, ok := aMap[i]; ok {
			delete(aMap, i)
		} else {
			bOnly = append(bOnly, i)
		}
	}
	for i := range aMap {
		aOnly = append(aOnly, i)
	}
	return aOnly, bOnly
}
