package user

func countSimilar(a, b []string) int {
	set := make(map[string]struct{})
	count := 0

	// put all items from a into a set
	for _, v := range a {
		set[v] = struct{}{}
	}

	// check each item in b
	for _, v := range b {
		if _, exists := set[v]; exists {
			count++
		}
	}

	return count
}
