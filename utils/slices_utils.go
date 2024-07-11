package utils

func RemoveElementsFromSlice[T comparable](slice []T, indexes []int) []T {
	// Create a map to store indexes to be removed
	indexMap := make(map[int]bool)
	for _, idx := range indexes {
		indexMap[idx] = true
	}

	// Create a new slice to store the result
	result := make([]T, 0, len(slice)-len(indexes))

	// Iterate through the original slice and add elements not in indexMap to result
	for idx, val := range slice {
		if !indexMap[idx] {
			result = append(result, val)
		}
	}

	return result
}

// SliceContains checks if a slice of any type contains a particular element
func SliceContains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
