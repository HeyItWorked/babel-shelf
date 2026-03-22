// models.gleam — Defines the Book type and BookStatus custom type.
// Shared across all layers. Equivalent of models.go and models.ts.
//
// Go uses string constants — typos only caught at runtime.
// TypeScript uses union types — caught at compile time.
// Gleam uses custom types — caught at compile time + exhaustive matching.

/// The three allowed reading statuses.
/// Unlike Go/TS where these are just strings, Gleam makes them
/// their own type — the compiler rejects anything else.
pub type BookStatus {
  WantToRead
  Reading
  Finished
}

/// A book on the shelf. Immutable — once created, never changed.
/// To "update" a field, you create a new Book with the changed value.
pub type Book {
  Book(id: Int, title: String, author: String, status: BookStatus)
}

/// Convert a BookStatus to the string the database stores.
/// Pattern matching guarantees every variant is handled —
/// if you add a new status, this won't compile until you add a case.
pub fn status_to_string(status: BookStatus) -> String {
  case status {
    WantToRead -> "want to read"
    Reading -> "reading"
    Finished -> "finished"
  }
}

/// Convert a database string back to a BookStatus.
/// Returns Result because the string might be invalid —
/// Go would return (status, error), TS would throw, Gleam returns Result.
pub fn status_from_string(s: String) -> Result(BookStatus, Nil) {
  case s {
    "want to read" -> Ok(WantToRead)
    "reading" -> Ok(Reading)
    "finished" -> Ok(Finished)
    _ -> Error(Nil)
  }
}
