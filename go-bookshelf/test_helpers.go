// test_helpers.go — Reusable functions for tests.
// Keeps test files focused on what's being tested, not boilerplate.
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// sendRequest builds and sends an HTTP request through the router,
// returns the recorder so you can check the response.
func sendRequest(method, url, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req, _ = http.NewRequest(method, url, strings.NewReader(body))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := setupRouter()
	handler.ServeHTTP(rr, req)
	return rr
}

// sendAndExpect sends a request and checks the status code for you.
// Returns the recorder so you can still check the body.
func sendAndExpect(t *testing.T, method, url, body string, expectedStatus int) *httptest.ResponseRecorder {
	rr := sendRequest(method, url, body)
	if rr.Code != expectedStatus {
		t.Errorf("got status %d, want %d", rr.Code, expectedStatus)
	}
	return rr
}

// decodeBook reads the response body into a single Book struct.
func decodeBook(t *testing.T, rr *httptest.ResponseRecorder) Book {
	var got Book
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(&got)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	return got
}

// decodeBooks reads the response body into a slice of Books.
func decodeBooks(t *testing.T, rr *httptest.ResponseRecorder) []Book {
	var got []Book
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(&got)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	return got
}
