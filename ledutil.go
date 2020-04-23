package lcv

import "log"

// Error checking function
func chk(err error) {
	if err != nil {
		log.Panic(err)
	}
}
