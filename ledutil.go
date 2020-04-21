package lcv

import "log"

// General Functions
func chk(err error) {
	if err != nil {
		log.Panic(err)
	}
}
