// handlers_test.go — Integration tests. Sends real HTTP requests through
// the full stack (router → handler → db → Postgres) and checks the responses.
// TDD: write tests first, watch them fail, then build to make them pass.
package main

import (
	"net/http"
	"testing"
)

func TestCreateBook(t *testing.T){
	// send POST /books with a new book
	rr := sendAndExpect(t, "POST", "/books", `{"title": "Dune", "author": "Frank Herbert"}`, http.StatusCreated)

	got := decodeBook(t, rr)

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
	// send GET /books
	rr := sendAndExpect(t, "GET", "/books", "", http.StatusOK)

	got := decodeBooks(t, rr)

	if len(got) == 0 {
		t.Errorf("expected at least one book, got empty list")
	}
}

func TestGetBook(t *testing.T) {
	// send GET /books/1
	rr := sendAndExpect(t, "GET", "/books/1", "", http.StatusOK)

	got := decodeBook(t, rr)

	if got.Id == 0 {
		t.Errorf("expected an id, got 0")
	}
	if got.Title == "" {
		t.Errorf("expected a title, got empty")
	}
	if got.Author == "" {
		t.Errorf("expected an author, got empty")
	}
	if got.Status == "" {
		t.Errorf("expected a status, got empty")
	}
}

func TestGetBookNotFound(t *testing.T) {
	// send GET /books/99999, an id that doesn't exist
	sendAndExpect(t, "GET", "/books/99999", "", http.StatusNotFound)
}

func TestUpdateBook(t *testing.T) {
	// send PUT /books/1 with JSON body containing new status
	rr := sendAndExpect(t, "PUT", "/books/1", `{"status": "reading"}`, http.StatusOK)

	got := decodeBook(t, rr)

	if got.Status != StatusReading {
		t.Errorf("got status %q, want %q", got.Status, StatusReading)
	}
}

func TestDeleteBook(t *testing.T) {
	// send DELETE /books/1
	sendAndExpect(t, "DELETE", "/books/1", "", http.StatusNoContent)
}

func TestCreateBookNoTitle(t *testing.T) {
	// send POST /books with JSON missing the title field
	sendAndExpect(t, "POST", "/books", `{"author": "Frank Herbert"}`, http.StatusBadRequest)
}

func TestCreateBookInvalidStatus(t *testing.T) {
	// send POST /books with a status that isn't one of the 3 allowed
	sendAndExpect(t, "POST", "/books", `{"title": "Dune", "author": "Frank Herbert", "status": "burned"}`, http.StatusBadRequest)
}
