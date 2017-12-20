package util

// The MapDiff file will hold all differences
// between two maps.
type MapDiff struct {
	Additions map[interface{}]interface{}
	Updates   map[interface{}]interface{}
	Deletions map[interface{}]interface{}
}

// DiffMaps takes in 2 maps with any key/value types & returns a map with the differences,
// additions, updates and deletions. This function assumes both maps are sorted. If something is
// in "a" and not "b" we class it as an addition, if it's in "b" but not "a" we class it as a
// deletion. We assume "a" is the primary map. This function only compares at the first level.
func DiffMapsKeys(a, b map[interface{}]interface{}) *MapDiff {
	var existsA bool
	var existsB bool

	// Somewhere to store our results
	result := &MapDiff{
		Additions: map[interface{}]interface{}{},
		Updates:   map[interface{}]interface{}{},
		Deletions: map[interface{}]interface{}{},
	}

	// Go over each item in a
	for i, data := range a {

		// Does this also exist in b?
		if _, existsB = b[i]; existsB {
			// Add to another list to compare after
			result.Updates[i] = []interface{}{
				data,
				b[i],
			}
			continue
		}

		// We didn't see it in b, so we'll mark it as an addition
		result.Additions[i] = data
	}

	// Go over each item in b
	for i, data := range b {

		// Does this also exist in a?
		if _, existsA = a[i]; existsA {
			// We will have already picked it up,
			// so we can just skip it.
			continue
		}

		// We didn't see it in a, so we'll mark it as deletion
		result.Deletions[i] = data
	}

	return result
}
