// handlers_test.go — Integration tests. Sends real HTTP requests through
// the full stack (router → handler → db → Postgres) and checks the responses.
// TDD: write tests first, watch them fail, then build to make them pass.
package main

import (
	"testing"
	"net/http"
	"strings"
	"net/http/httptest"
	"encoding/json"
)

func TestCreateBook(t *testing.T){
	// build the request — assemble the envelope
	body := strings.NewReader(`{"title": "Dune", "author": "Frank Herbert"}`)
	req, err := http.NewRequest("POST", "/books", body)
	if err != nil{
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request — handler writes its response into the recorder
	responseRecorder := httptest.NewRecorder()
	handler := setupRouter()
	handler.ServeHTTP(responseRecorder, req)

	// check the result — read the notebook
	if responseRecorder.Code != http.StatusCreated {
		t.Errorf("got status %d, want %d", responseRecorder.Code, http.StatusCreated)
	}

	var got Book
	decoder := json.NewDecoder(responseRecorder.Body)
	err = decoder.Decode(&got)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if got.Title != "Dune" {
		t.Errorf("got title %q, want %q", got.Title, "Dune")
	}
	if got.Author != "Frank Herbert" {
		t.Errorf("got author %q, want %q", got.Author, "Frank Herbert")
	}
	if got.Status != StatusWantToRead {
		t.Errorf("got status %q, want %q", got.Status, StatusWantToRead)
	}
	if got.Id == 0 {
		t.Errorf("expected an id, got 0")
	}
}

func TestGetBooks(t *testing.T) {
	// build the request — GET /books, no body needed
	// send the request — handler writes list of books into the recorder
	// check the result — decode into []Book (a slice, not a single book)
	// verify the list is not empty and contains books we expect
}

func TestGetBook(t *testing.T) {
	// build the request — GET /books/1, no body needed
	// send the request — handler writes one book into the recorder
	// check the result — decode into a single Book
	// verify the id, title, author, status match what we created
}

func TestGetBookNotFound(t *testing.T) {
	// build the request — GET /books/99999, an id that doesn't exist
	// send the request — handler should respond with not found
	// check the result — expect 404 status code
}

func TestUpdateBook(t *testing.T) {
	// build the request — PUT /books/1 with JSON body containing new status
	// send the request — handler updates the book and writes it into the recorder
	// check the result — decode into a Book
	// verify the updated fields changed and the rest stayed the same
}

func TestDeleteBook(t *testing.T) {
	// build the request — DELETE /books/1, no body needed
	// send the request — handler deletes the book
	// check the result — expect 200 or 204 (no content)
	// optionally: GET /books/1 again and expect 404 to confirm it's gone
}

func TestCreateBookNoTitle(t *testing.T) {
	// build the request — POST /books with JSON missing the title field
	// send the request — handler should reject it
	// check the result — expect 400 (bad request)
}

func TestCreateBookInvalidStatus(t *testing.T) {
	// build the request — POST /books with a status that isn't one of the 3 allowed
	// send the request — handler should reject it
	// check the result — expect 400 (bad request)
}