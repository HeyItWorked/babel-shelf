// handlers.go — Logic layer. One function per CRUD operation.
// Receives HTTP requests, talks to the DB layer, sends responses.
//
// Pattern: read from user → talk to DB → write back to user
//
// Rule: every path through a handler must write to w exactly once.
// Error path? → http.Error(w, ...) + return
// Happy path? → respondJSON(w, ...) or w.WriteHeader(...)
// If you hit the closing } without writing to w, something's wrong.
//
//   Handler      | Input         | DB call      | Response
//   -------------|---------------|--------------|----------
//   CreateBook   | JSON body     | insertBook   | 201 + book
//   GetBooks     | nothing       | getAllBooks  | 200 + list
//   GetBook      | id from URL   | getBookById  | 200 + book
//   UpdateBook   | id + JSON     | updateBook   | 200 + book
//   DeleteBook   | id from URL   | deleteBook   | 204
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// parseId — extracts the id from the URL path and converts to int
func parseId(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

// respondJSON — writes a status code and JSON body to the response
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// CreateBook — POST /books
// Reads JSON body → calls insertBook → responds with the new book
func CreateBook(w http.ResponseWriter, r *http.Request) {
	// read from user — decode JSON body into a Book
	var book Book
	// Decode over Unmarshal — reads from stream directly, standard choice for HTTP bodies
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// validate — title and author are required
	if book.Title == "" || book.Author == "" {
		http.Error(w, "title and author are required", http.StatusBadRequest)
		return
	}

	// validate — status must be one of the 3 allowed values, default if empty
	switch book.Status {
	case "":
		book.Status = StatusWantToRead
	case StatusWantToRead, StatusReading, StatusFinished:
		// valid, do nothing
	default:
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	// talk to DB — insert the book
	book, err = insertBook(book)
	if err != nil {
		http.Error(w, "could not create book", http.StatusInternalServerError)
		return
	}

	// write back to user
	respondJSON(w, http.StatusCreated, book)
}

// GetBooks — GET /books
// No input needed → calls getAllBooks → responds with the list
func GetBooks(w http.ResponseWriter, r *http.Request) {
	// no input to read — just talk to DB
	books, err := getAllBooks()
	if err != nil {
		http.Error(w, "could not get books", http.StatusInternalServerError)
		return
	}

	// write back to user
	respondJSON(w, http.StatusOK, books)
}

// GetBook — GET /books/{id}
// Reads id from URL → calls getBookById → responds with the book
func GetBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	book, err := getBookById(id)
	if err != nil {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}
	respondJSON(w, http.StatusOK, book)
}

// UpdateBook — PUT /books/{id}
// Reads id from URL + JSON body → calls updateBook → responds with updated book
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var book Book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	book, err = updateBook(id, book)
	if err != nil {
		http.Error(w, "can't update book", http.StatusBadRequest)
		return
	}
	respondJSON(w, http.StatusOK, book)
}

// DeleteBook — DELETE /books/{id}
// Reads id from URL → calls deleteBook → responds with no content
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = deleteBook(id)
	if err != nil {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	// 204 = "it worked, nothing to show you" — sends the status code only, no JSON body.
	// every HTTP response has a status code; this one just has no body attached.
	w.WriteHeader(http.StatusNoContent)
}
