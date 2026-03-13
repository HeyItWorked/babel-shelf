// models.go — Defines the Book struct. This is what a book looks like in Go.
// Shared across all layers.
package main

// Go doesn't have enums. The convention is to use constants instead.
// These are just named strings — the compiler catches typos if you
// use StatusReading instead of the raw string "reading".
// The actual constraint is enforced by the database (CHECK in init.sql).
const (
	StatusWantToRead = "want to read"
	StatusReading    = "reading"
	StatusFinished   = "finished"
)

type Book struct {
    Id     int     `json:"id"`
    Title  string `json:"title"`
	Author string `json:"author"`
	Status string  `json:"status"`
}
