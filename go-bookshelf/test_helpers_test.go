// test_helpers.go — Reusable functions for tests.
// Keeps test files focused on what's being tested, not boilerplate.
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

// ERROR: tests panic with nil pointer dereference on db.QueryRow
// EXPLANATION: tests don't call main(), so db is never set — it's nil
// FIX: TestMain runs before all tests and connects to the database
// NOTE: this file must end in _test.go — Go only runs TestMain from test files
func TestMain(m *testing.M) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://shelf:shelf@db:5432/bookshelf?sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		fmt.Printf("could not connect to test database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	os.Exit(m.Run()) // runs all tests, exits with their result code
}

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
