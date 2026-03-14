// db.ts — Translator layer. Speaks SQL to Postgres and converts
// the results into TypeScript objects. All database queries live here.
// Equivalent of db.go in Go.

import {Pool} from "pg"
import {Book} from './models'

// shared database connection — equivalent of Go's `var db *sql.DB`
export let pool: Pool

// lets index.ts and tests assign the pool from outside
// can't just export and reassign — ES modules give you a copy, not the original.
// setter writes to the actual variable inside db.ts.
export function setPool(p: Pool) { pool = p }

// pg auto-maps columns to object keys — no manual Scan like Go
export async function insertBook(book: Book): Promise<Book> {
    const sql = `INSERT INTO books (title, author, status) VALUES ($1, $2, $3) RETURNING id, title, author, status`
    const result = await pool.query<Book>(sql, [book.title, book.author, book.status])
    return result.rows[0]
}

// no defer rows.Close() needed — pg handles cleanup automatically
export async function getAllBooks(): Promise<Book[]> {
    const sql = `SELECT id, title, author, status FROM books`
    const results = await pool.query<Book>(sql)
    return results.rows
}

export async function getBookById(id: number): Promise<Book | undefined> {
    const sql = `SELECT id, title, author, status FROM books WHERE id = $1`
    const result = await pool.query<Book>(sql, [id])
    return result.rows[0]
}

// same pattern as insertBook — RETURNING gives back the row in one trip
export async function updateBook(id: number, book: Book): Promise<Book | undefined> {
    const sql = `UPDATE books SET title=$1, author=$2, status=$3 WHERE id=$4 RETURNING id, title, author, status`
    const result = await pool.query<Book>(sql, [book.title, book.author, book.status, id])
    return result.rows[0]
}

// no RETURNING — nothing comes back, just check rowCount (like Go's RowsAffected)
export async function deleteBook(id: number): Promise<boolean> {
    const sql = `DELETE FROM books WHERE id = $1`
    const result = await pool.query(sql, [id])
    if (result.rowCount === null || result.rowCount === 0) {
        return false
    }
    return true
}
