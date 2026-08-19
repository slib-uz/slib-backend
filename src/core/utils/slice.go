package utils

// In checks if value exists in slice
func In[T comparable](value T, arr []T) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
