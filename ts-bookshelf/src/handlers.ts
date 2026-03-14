// handlers.ts — Logic layer. One function per CRUD operation.
// Receives HTTP requests, talks to the DB layer, sends responses.
// Equivalent of handlers.go in Go.
//
// Pattern: read from user → validate → talk to DB → write back to user
//
// Rule: every handler must return exactly one response.
// Error path? → return c.text("error", statusCode)
// Happy path? → return c.json(data, statusCode)
//
//   Handler      | Input         | DB call       | Response
//   -------------|---------------|---------------|----------
//   createBook   | JSON body     | insertBook    | 201 + book
//   getBooks     | nothing       | getAllBooks    | 200 + list
//   getBook      | id from URL   | getBookById   | 200 + book
//   updateBook   | id + JSON     | updateBook    | 200 + book
//   deleteBook   | id from URL   | deleteBook    | 204

import type { Context } from "hono"
import { insertBook, getAllBooks, getBookById, updateBook, deleteBook } from "./db"
import { VALID_STATUSES, STATUS_WANT_TO_READ } from "./models"

// ── createBook ── POST /books
// read JSON body → validate title+author → validate status → call insertBook → 201
export async function createBook(c: Context) {
    // const body = await c.req.json()
    // validate: if (!body.title || !body.author) → 400
    // validate: status must be in VALID_STATUSES, default to STATUS_WANT_TO_READ
    // const book = await insertBook(body)
    // return c.json(book, 201)
}

// ── getBooks ── GET /books
// no input → call getAllBooks → 200
export async function getBooks(c: Context) {
    // const books = await getAllBooks()
    // return c.json(books, 200)
}

// ── getBook ── GET /books/:id
// parse id from URL → call getBookById → 200 or 404
export async function getBook(c: Context) {
    // const id = Number(c.req.param("id"))
    // if (isNaN(id)) → 400
    // const book = await getBookById(id)
    // if (!book) → 404
    // return c.json(book, 200)
}

// ── updateBook ── PUT /books/:id
// parse id + JSON body → call updateBook → 200 or 400
export async function handleUpdateBook(c: Context) {
    // const id = Number(c.req.param("id"))
    // if (isNaN(id)) → 400
    // const body = await c.req.json()
    // const book = await updateBook(id, body)
    // if (!book) → 400
    // return c.json(book, 200)
}

// ── deleteBook ── DELETE /books/:id
// parse id from URL → call deleteBook → 204 or 404
export async function handleDeleteBook(c: Context) {
    // const id = Number(c.req.param("id"))
    // if (isNaN(id)) → 400
    // const deleted = await deleteBook(id)
    // if (!deleted) → 404
    // return c.body(null, 204)
}
