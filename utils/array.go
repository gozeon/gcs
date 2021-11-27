package utils

// Find @link https://yourbasic.org/golang/find-search-contains-slice/
// Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// @link https://stackoverflow.com/questions/15323767/does-go-have-if-x-in-construct-similar-to-python#answer-34599643
// func Search()
