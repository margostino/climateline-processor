package common

import "log"

func Check(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
