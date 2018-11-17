package Linearsearch

// Contains tells whether a contains x.
/*func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}*/

func Contains(a []rune, x rune) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}