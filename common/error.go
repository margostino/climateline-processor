package common

import "log"

func Check(err error, message string) {
	if err != nil {
		log.Fatalf("Error: %s - %s\n", err.Error(), message)
	}
}

func SilentCheck(err error, message string) {
	if err != nil {
		log.Printf("Error: %s - %s\n", err.Error(), message)
	}
}

func IsError(err error, message string) bool {
	if err != nil {
		log.Printf("Error: %s - %s\n", err.Error(), message)
		return true
	}
	return false
}
