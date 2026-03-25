// db.gleam — Database layer. Speaks SQL to Postgres and converts
// the results into Gleam types. All database queries live here.
// Equivalent of db.go and db.ts in the other implementations.
//
// Key differences from Go/TS:
//   - No global mutable state — the connection is passed as an argument.
//   - Row decoder converts (id, title, author, status) tuples to Book records.
//   - `use` keyword chains fallible operations (replaces Go's `if err != nil`).

import gleam/dynamic/decode
import gleam/result
import gleam/string
import models.{
  type Book, type BookStatus, Book, status_from_string, status_to_string,
}
import pog

/// Application-level errors returned by DB functions.
/// NotFound means 0 rows matched — used for get-by-id and delete.
/// DatabaseError wraps pog.QueryError as a debug string.
pub type AppError {
  NotFound
  DatabaseError(String)
}

/// Decode a (id, title, author, status) row into a Book.
/// The `let assert` is safe because the CHECK constraint in init.sql
/// guarantees status is always one of the three valid values.
pub fn book_decoder() -> decode.Decoder(Book) {
  use id <- decode.field(0, decode.int)
  use title <- decode.field(1, decode.string)
  use author <- decode.field(2, decode.string)
  use status_str <- decode.field(3, decode.string)
  let assert Ok(status) = status_from_string(status_str)
  decode.success(Book(id:, title:, author:, status:))
}

/// Convert a pog.QueryError into an AppError.
pub fn map_error(res: Result(a, pog.QueryError)) -> Result(a, AppError) {
  result.map_error(res, fn(e) { DatabaseError(string.inspect(e)) })
}

/// INSERT a new book, return it with the generated id.
/// Same INSERT ... RETURNING pattern as Go and TS — one round trip.
pub fn insert_book(
  conn: pog.Connection,
  title: String,
  author: String,
  status: BookStatus,
) -> Result(Book, AppError) {
  let sql =
    "INSERT INTO books (title, author, status) VALUES ($1, $2, $3)
     RETURNING id, title, author, status"

  use returned <- result.try(
    pog.query(sql)
    |> pog.parameter(pog.text(title))
    |> pog.parameter(pog.text(author))
    |> pog.parameter(pog.text(status_to_string(status)))
    |> pog.returning(book_decoder())
    |> pog.execute(conn)
    |> map_error,
  )

  case returned.rows {
    [book] -> Ok(book)
    _ -> Error(DatabaseError("insert returned no rows"))
  }
}
