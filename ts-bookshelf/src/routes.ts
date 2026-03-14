// routes.ts — Route mapping. Which URL goes to which handler.
// Equivalent of setupRouter() in Go's main.go.
//
// Separated from index.ts so tests can import the app
// without booting the server.

import { Hono } from "hono"
import { createBook, getBooks, getBook, handleUpdateBook, handleDeleteBook } from "./handlers"

const app = new Hono()

// POST   /books      → createBook
// GET    /books      → getBooks
// GET    /books/:id  → getBook       (:id not {id} like Go)
// PUT    /books/:id  → handleUpdateBook
// DELETE /books/:id  → handleDeleteBook

export default app
