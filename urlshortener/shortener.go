package urlshortener

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/margostino/climateline-processor/common"
	"net/http"
	"os"
)

var bitlyDomain = "bit.ly"
var htpClient = &http.Client{}

var token = os.Getenv("BITLY_TOKEN")
var endpoint = os.Getenv("BITLY_SHORTENER_ENDPOINT")

func Shorten(url string) string {
	var urlShortenerResponse map[string]interface{}
	request := getRequest(url)
	response, err := htpClient.Do(request)
	if !common.IsError(err, "when shorting url") && (response.StatusCode == 201 || response.StatusCode == 200) {
		err = json.NewDecoder(response.Body).Decode(&urlShortenerResponse)
		if !common.IsError(err, "when decoding shortener url response") {
			return urlShortenerResponse["link"].(string)
		}
	}
	return ""
}

func getRequest(url string) *http.Request {
	jsonRequest, err := json.Marshal(map[string]string{"long_url": url})
	if !common.IsError(err, "when marshalling request") {
		request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonRequest))
		if !common.IsError(err, "when creating URL shortener request") {
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			request.Header.Set("Content-Type", "application/json")
			return request
		}
	}
	return nil
}
