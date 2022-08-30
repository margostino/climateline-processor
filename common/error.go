package common

import "log"

func Check(err error, message string) {
	if err != nil {
		log.Fatalf("Error: %s - %s", err.Error(), message)
	}
}
