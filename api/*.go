package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func S(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Hello World from Go"
	resp["language"] = "goaaaaa"
	resp["cloud"] = "Hosted on Vercel!"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("Error happened in JSON marshal. Err: %s\n", err)
	} else {
		w.Write(jsonResp)
	}
	return
}
