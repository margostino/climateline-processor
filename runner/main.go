package main

import (
	"github.com/margostino/climateline-processor/api"
	"log"
	"net/http"
)

func main() {
	//http.HandleFunc("/", api.RunJob)
	http.HandleFunc("/bot", api.Bot)
	http.HandleFunc("/publisher-job", api.ExecutePublisherJob)
	http.HandleFunc("/collector-job", api.ExecuteCollectorJob)
	http.HandleFunc("/news", api.News)
	http.HandleFunc("/cache", api.Cache)
	log.Println("Starting Server in :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
