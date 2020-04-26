package lcv

import "log"

// Checks the error and panics if one occurred
func chk(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// Returns the index of a string in a string slice, -1 returned if not found
func stringpos(s []string, value string) int {
	for p, v := range s {
		if v == value {
			return p
		}
	}
	return -1
}

// Returns the index of an int in a string slice, -1 returned if not found
func intpos(s []int, value int) int {
	for p, v := range s {
		if v == value {
			return p
		}
	}
	return -1
}
