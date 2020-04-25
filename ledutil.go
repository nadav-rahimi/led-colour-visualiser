package lcv

import "log"

// Error checking function
func chk(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// Slice functions
func pos(s []string, value string) int {
	for p, v := range s {
		if v == value {
			return p
		}
	}
	return -1
}
