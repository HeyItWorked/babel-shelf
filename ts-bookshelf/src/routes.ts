// routes.ts — Route mapping. Which URL goes to which handler.
// Separated from index.ts so tests can import the app without booting the server.

import { Hono } from "hono"
import { createBook, getBooks, getBook, handleUpdateBook, handleDeleteBook } from "./handlers"

const app = new Hono()

// /books     → create, list
// /books/:id → get, update, delete
app.post("/books", createBook)
app.get("/books", getBooks)
app.get("/books/:id", getBook)
app.put("/books/:id", handleUpdateBook)
app.delete("/books/:id", handleDeleteBook)

export default app
