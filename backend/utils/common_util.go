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

func StringXOR(list1, list2 []string) []string {
    set1 := make(map[string]bool)
    for _, s := range list1 {
        set1[s] = true
    }
    set2 := make(map[string]bool)
    for _, s := range list2 {
        set2[s] = true
    }

    var c []string
    for _, s := range list1 {
        if !set2[s] {
          c = append(c, s)
        }
    }
    for _, s := range list2 {
        if !set1[s] {
          c = append(c, s)
        }
    }
    return c
}
