// handlers.ts — Logic layer. read → validate → db → respond.
// Context (c) = request + response in one object (Go splits these into r and w)
// c.req.json() / c.req.param("id") to read, c.json() / c.text() to respond

import type { Context } from "hono"
import { insertBook, getAllBooks, getBookById, updateBook, deleteBook } from "./db"
import { VALID_STATUSES, STATUS_WANT_TO_READ } from "./models"

// POST /books — body → validate title+author → validate status → insertBook → 201
export async function createBook(c: Context): Promise<Response> {
    const body = await c.req.json()
    if (!body.title || !body.author) {
        return c.text("title and author are required", 400)
    }
    if (!body.status) {
        body.status = STATUS_WANT_TO_READ
    } else if (!VALID_STATUSES.includes(body.status)) {
        return c.text("invalid status", 400)
    }
    const book = await insertBook(body)
    return c.json(book, 201)
}

// GET /books — getAllBooks → 200
// pure wrapper
export async function getBooks(c: Context): Promise<Response> {
    const books = await getAllBooks()
    return c.json(books, 200)
}

// GET /books/:id — parse id → getBookById → 200 or 404
export async function getBook(c: Context): Promise<Response> {
    const raw = c.req.param("id")   // URL param is always a string
    const id = Number(raw)           // convert to number
    if (isNaN(id)) {
        return c.text("invalid id", 400)
    }
    const book = await getBookById(id)
    if (!book) {
        return c.text("book not found", 404)
    }
    return c.json(book, 200)
}

// PUT /books/:id — parse id + body → updateBook → 200 or 400
export async function handleUpdateBook(c: Context): Promise<Response> {
    const raw = c.req.param("id")
    const id = Number(raw)
    if (isNaN(id)) {
        return c.text("invalid id", 400)
    }
    const body = await c.req.json()
    const book = await updateBook(id, body)
    if (!book) {
        return c.text("can't update book", 400)
    }
    return c.json(book, 200)
}

// DELETE /books/:id — parse id → deleteBook → 204 or 404
export async function handleDeleteBook(c: Context): Promise<Response> {
    const raw = c.req.param("id")
    const id = Number(raw)
    if (isNaN(id)) {
        return c.text("invalid id", 400)
    }
    const deleted = await deleteBook(id)
    if (!deleted) {
        return c.text("book not found", 404)
    }
    return c.body(null, 204)
}
