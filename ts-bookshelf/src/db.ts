// db.ts — Translator layer. Speaks SQL to Postgres and converts
// the results into TypeScript objects. All database queries live here.
// Equivalent of db.go in Go.

// ── imports ──
// Pool from pg — manages database connections
// Book from models — the type we return

import {Pool} from "pg"
import {Book} from './models'

// ── pool ──
// shared database connection — equivalent of Go's `var db *sql.DB`
// why export? TS files are isolated — index.ts imports pool to set it up,
// after that only db.ts touches it.

export let pool: Pool

// ── insertBook ──
// INSERT a new book → returns it with the generated id
// pg auto-maps columns to object keys — no manual Scan like Go
export async function insertBook(book: Book): Promise<Book> {
    const sql: string = `INSERT INTO books (title, author, status) VALUES ($1, $2, $3) RETURNING id, title, author, status`
    // tells pg the rows coming back are Books objects
    const result = await pool.query<Book>(sql, [book.title, book.author, book.status])
    return result.rows[0]
}

// ── getAllBooks ──
// SELECT all books → returns the full array
// no defer rows.Close() needed — pg handles cleanup automatically
export async function getAllBooks(): Promise<Book[]> {
    const sql: string = `SELECT id, title, author, status FROM books`
    // each row is a Book
    const results = await pool.query<Book>(sql)
    return results.rows
}

// ── getBookById ──
// SELECT one book by id → returns the book, or undefined if not found
export async function getBookById(id: number): Promise<Book | undefined> {
    const sql = `SELECT id, title, author, status FROM books WHERE id = $1`
    const result = await pool.query<Book>(sql, [id])
    return result.rows[0]
}

// ── updateBook ──
// UPDATE a book's fields by id → returns updated book, or undefined if not found
// same pattern as insertBook — RETURNING gives back the row in one trip
export async function updateBook(id: number, book: Book): Promise<Book | undefined> {
    // SQL: UPDATE books SET title=$1, author=$2, status=$3 WHERE id=$4 RETURNING id, title, author, status
    // pool.query with [book.title, book.author, book.status, id]
    // return result.rows[0]
}

// ── deleteBook ──
// DELETE a book by id → returns true if deleted, false if it didn't exist
// no RETURNING — nothing comes back, just check rowCount (like Go's RowsAffected)
export async function deleteBook(id: number): Promise<boolean> {
    const sql = `DELETE FROM books WHERE id = $1`
    const result = await pool.query(sql, [id])
    if (result.rowCount === null || result.rowCount === 0) {
        return false
    }
    return true
}