package utils

func FilterStringArray(array []string, element string) []string {
	result := []string{}

	for _, a := range array {
		if a != element {
			result = append(result, a)
		}
	}

	return result
}
