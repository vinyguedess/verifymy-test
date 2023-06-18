package utils

// SliceContains checks if a slice contains an element
func SliceContains(slice []interface{}, element interface{}) bool {
	for _, sliceElement := range slice {
		if sliceElement == element {
			return true
		}
	}

	return false
}
