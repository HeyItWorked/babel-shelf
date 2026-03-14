// handlers.test.ts — Integration tests. Same 8 tests as Go.
// Sends requests through Hono → handler → db → Postgres.

import { describe, test, expect } from "bun:test"
import "./helpers"
import { sendAndExpect } from "./helpers"

describe("books", () => {

    // POST /books — create a book, check fields + id assigned
    test("create book", async () => {
        const res = await sendAndExpect("POST", "/books", { title: "Dune", author: "Frank Herbert"}, 201)
        // response object has json() built in
        const body = await res.json()

        expect(body.title).toBe("Dune")
        expect(body.author).toBe("Frank Herbert")
        expect(body.status).toBe("want to read")
        expect(body.id).toBeGreaterThan(0)
    })

    // GET /books — list all, check non-empty
    test("get books", async () => {
        // parse body, expect array length > 0
        const res = await sendAndExpect("GET", "/books", undefined, 200)
        const body = await res.json()

        expect(body.length).toBeGreaterThan(0)
    })

    // GET /books/1 — get one, check fields exist
    test("get book", async () => {
        // sendAndExpect("GET", "/books/1", undefined, 200)
        // parse body, expect id, title, author, status all present
        const res = await sendAndExpect("GET", "/books/1", undefined, 200)
        const body = await res.json()

        expect(body.id).toBeGreaterThan(0)
        expect(body.title).not.toBeEmpty()
        expect(body.author).not.toBeEmpty()
        expect(body.status).not.toBeEmpty()
    })

    // GET /books/99999 — not found
    test("get book not found", async () => {
        // no body (ba dum ts) found, just let SAE check status code
        await sendAndExpect("GET", "/books/99999", undefined, 404)
    })

    // PUT /books/1 — update status, check it changed
    test("update book", async () => {
        const res = await sendAndExpect("PUT", "/books/1", { title: "Dune", author: "Frank Herbert", status: "reading" }, 200)
        const body = await res.json()

        expect(body.status).toBe("reading")
    })

    // DELETE /books/1 — delete, check 204
    test("delete book", async () => {
        await sendAndExpect("DELETE", "/books/1", undefined, 204)
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
