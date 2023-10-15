package util

// Delete removes the elements s[i:j] from s, returning the modified slice.
// Delete panics if s[i:j] is not a valid slice of s.
// Delete modifies the contents of the slice s; it does not create a new slice.
// Delete is O(len(s)-(j-i)), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
//
// DISCLAIMER: This func is yoinked from golang experimental branch:
// https://cs.opensource.google/go/x/exp/+/0b5c67f0:slices/slices.go;l=156
func DelFromSlice[S ~[]E, E any](s S, i, j int) S {
	return append(s[:i], s[j:]...)
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
// DISCLAIMER: This func is yoinked from golang experimental branch:
// https://cs.opensource.google/go/x/exp/+/92128663:slices/slices.go;l=93
func SliceIndex[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// Contains reports whether v is present in s.
// DISCLAIMER: This func is yoinked from golang experimental branch:
// https://cs.opensource.google/go/x/exp/+/92128663:slices/slices.go;l=93
func SliceContains[S ~[]E, E comparable](s S, v E) bool {
	return SliceIndex(s, v) >= 0
}
