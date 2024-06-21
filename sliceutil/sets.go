package sliceutil

// Difference returns a new slice with elements that are in "set" but not in
// "others".
func Difference[T comparable](set []T, others ...[]T) []T {
	out := []T{}

	for _, item := range set {
		var found bool
		for _, o := range others {
			if Contains(o, item) {
				found = true
				break
			}
		}

		if !found {
			out = append(out, item)
		}
	}

	return out
}

// Complement returns the complement of the two lists; that is, the first return
// value will contain elements that are only in "a" (and not in "b"), and the
// second return value will contain elements that are only in "b" (and not in
// "a").
func Complement[T comparable](a, b []T) (aOnly, bOnly []T) {
	for _, i := range a {
		if !Contains(b, i) {
			aOnly = append(aOnly, i)
		}
	}

	for _, i := range b {
		if !Contains(a, i) {
			bOnly = append(bOnly, i)
		}
	}

	return aOnly, bOnly
}

// Intersection returns the elements common to both a and b
func Intersection[T comparable](a, b []T) []T {
	inter := []T{}
	hash := make(map[T]bool, len(a))
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

// IntersectionOfMany returns the elements common to all arguments
func IntersectionOfMany[T comparable](slices ...[]T) []T {
	if len(slices) == 1 {
		return slices[0]
	}
	common := slices[0]
	for i, s := range slices {
		if i == 0 {
			continue
		}
		common = Intersection(common, s)
		if len(common) == 0 {
			break
		}
	}
	return common
}
