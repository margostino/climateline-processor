package main

import (
	"github.com/margostino/climateline-processor/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", api.RunJob)
	log.Println("Starting Server in :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
