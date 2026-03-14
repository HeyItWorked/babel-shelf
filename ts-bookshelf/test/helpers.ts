// helpers.ts — Reusable functions for tests.
// Go uses TestMain, Bun uses beforeAll/afterAll — same idea.

import { beforeAll, afterAll, expect } from "bun:test"
import { Pool } from "pg"
import { pool, setPool } from "../src/db"
import app from "../src/routes"

// ── setup — connect to Postgres before tests ──
beforeAll(async() => {
    // connectionString is shortcut for { host, port, user, ... } — matches DATABASE_URL format
    const databaseUrl = process.env.DATABASE_URL || "postgres://shelf:shelf@db:5432/bookshelf?sslmode=disable"
    const p = new Pool({ connectionString: databaseUrl })
    // can't import pool directly and reassign — ES modules give you a copy, not the original
    setPool(p)
})

// ── teardown — close pool after tests ──
afterAll(async() => {
    await pool.end()
})

// ── sendRequest — wraps app.request(), returns Response ──
export async function sendRequest(method: string, url: string, body?: object): Promise<Response> {
    return await app.request(url, {
        method: method,
        body: body ? JSON.stringify(body) : undefined,
        headers: {"Content-Type": "application/json"}
    })
}

// ── sendAndExpect — sendRequest + assert status code ──
export async function sendAndExpect(method: string, url: string, body: object | undefined, expectedStatus: number): Promise<Response>{
    const result = await sendRequest(method, url, body)
    expect(result.status).toBe(expectedStatus)
    return result
}
