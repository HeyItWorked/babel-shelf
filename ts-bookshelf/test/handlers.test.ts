// handlers.test.ts — Integration tests. Same 8 tests as Go.
// Sends requests through Hono → handler → db → Postgres.
//
// Go tests hardcoded /books/1 — worked because Go ran on a fresh DB.
// Shared DB means IDs are unpredictable, so we capture the id from create
// and use it in get/update/delete.

import { describe, test, expect } from "bun:test"
import "./helpers"
import { sendAndExpect } from "./helpers"

// shared across tests — set by create, used by get/update/delete
let bookId: number

describe("books", () => {

    // POST /books — create a book, capture its id for later tests
    test("create book", async () => {
        const res = await sendAndExpect("POST", "/books", { title: "Dune", author: "Frank Herbert"}, 201)
        const body = await res.json()

        expect(body.title).toBe("Dune")
        expect(body.author).toBe("Frank Herbert")
        expect(body.status).toBe("want to read")
        expect(body.id).toBeGreaterThan(0)

        bookId = body.id
    })

    // GET /books — list all, check non-empty
    test("get books", async () => {
        const res = await sendAndExpect("GET", "/books", undefined, 200)
        const body = await res.json()

        expect(body.length).toBeGreaterThan(0)
    })

    // GET /books/:id — get the book we created
    test("get book", async () => {
        const res = await sendAndExpect("GET", `/books/${bookId}`, undefined, 200)
        const body = await res.json()

        expect(body.id).toBe(bookId)
        expect(body.title).not.toBeEmpty()
        expect(body.author).not.toBeEmpty()
        expect(body.status).not.toBeEmpty()
    })

    // GET /books/99999 — not found
    test("get book not found", async () => {
        await sendAndExpect("GET", "/books/99999", undefined, 404)
    })

    // PUT /books/:id — update the book we created
    test("update book", async () => {
        const res = await sendAndExpect("PUT", `/books/${bookId}`, { title: "Dune", author: "Frank Herbert", status: "reading" }, 200)
        const body = await res.json()

        expect(body.status).toBe("reading")
    })

    // DELETE /books/:id — delete the book we created
    test("delete book", async () => {
        await sendAndExpect("DELETE", `/books/${bookId}`, undefined, 204)
    })

    // POST /books — missing title, check 400
    test("create book no title", async () => {
        await sendAndExpect("POST", "/books", { author: "Frank Herbert" }, 400)
    })

    // POST /books — invalid status, check 400
    test("create book invalid status", async () => {
        await sendAndExpect("POST", "/books", { title: "Dune", author: "Frank Herbert", status: "burned" }, 400)
    })
})
