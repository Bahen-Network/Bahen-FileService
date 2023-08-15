package util

import "log"

func HandleErr(err error, msg string) bool {
	if err != nil {
		log.Printf("%s - error: %v", msg, err)
		return true
	}
	return false
}
