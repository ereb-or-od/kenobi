package utilities

func CompareSlices(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func ContainsInStringSlice(slice []string, value string) bool{
	if slice != nil && len(slice) > 0{
		for _, element := range slice{
			if element == value{
				return true
			}
		}
	}
	return false
}