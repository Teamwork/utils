package mathutil

// ComplementsInt takes two int64 slices, and returns the complements of the two
// lists removing repeated values.
func ComplementsInt(a, b []int64) (aOnly, bOnly []int64) {
	aMap := make(map[int64]struct{}, len(a))
	for _, i := range a {
		aMap[i] = struct{}{}
	}
	bMap := make(map[int64]struct{}, len(b))
	for _, i := range b {
		if _, ok := aMap[i]; ok {
			delete(aMap, i)
		} else {
			bMap[i] = struct{}{}
		}
	}
	for i := range aMap {
		aOnly = append(aOnly, i)
	}
	for i := range bMap {
		bOnly = append(bOnly, i)
	}
	return aOnly, bOnly
}
