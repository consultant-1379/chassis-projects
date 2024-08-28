package slicetool

// UniqueInt remove duplicate values from slice []int
func UniqueInt(intSlice []int) []int {
	if intSlice == nil {
		return nil
	}
	keys := make(map[int]bool)
	var list []int
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// UniqueString remove duplicate values from slice []string
func UniqueString(stringSlice []string) []string {
	if stringSlice == nil {
		return nil
	}
	keys := make(map[string]bool)
	var list []string
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
