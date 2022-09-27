package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func assertError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func assertValidInput(isValidMessage bool, input string, t *testing.T) {
	if !isValidMessage {
		t.Errorf("input [ %s ] is not valid", input)
	}
}

func assertStatus(response *httptest.ResponseRecorder, expectedStatus int, t *testing.T) {
	if status := response.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func assertResponse(unexpectedCondition bool, current interface{}, expected interface{}, t *testing.T) {
	if unexpectedCondition {
		t.Errorf("handler returned unexpected response: got %v want %v", current, expected)
	}
}
