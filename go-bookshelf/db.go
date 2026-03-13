// db.go — Translator layer. Speaks SQL to Postgres and converts
// the results into Go structs. All database queries live here.
package main

import "fmt"

// insertBook — INSERT a new book, return it with the new id.
//
// How book flows through QueryRow and Scan:
//   book  ──read──→  QueryRow  ──writes to──→  Postgres
//   book  ←──write──  Scan     ←──reads from──  Postgres
func insertBook(book Book) (Book, error) {
	// RETURNING is Postgres-specific — gives back the inserted row in one trip
	row := db.QueryRow(
		`INSERT INTO books (title, author, status) VALUES ($1, $2, $3)
		 RETURNING id, title, author, status`,
		book.Title, book.Author, book.Status,
	)

	err := row.Scan(&book.Id, &book.Title, &book.Author, &book.Status)
	if err != nil {
		return book, err
	}
	return book, nil
}

// getAllBooks — SELECT all books.
//
// Why defer rows.Close() here but not in insertBook?
//   QueryRow returns one row — cleans up automatically after Scan.
//   Query holds an open connection until you close it.
//   defer rows.Close() makes sure we hang up, even if we error mid-loop.
func getAllBooks() ([]Book, error) {
	// explicit columns instead of SELECT * — Scan expects exactly 4 fields,
	// so SELECT * would break if a column is added later
	rows, err := db.Query("SELECT id, title, author, status FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Status)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// getBookById — SELECT one book by id.
// Returns Book{} (zero value) on error since structs can't be nil.
func getBookById(id int) (Book, error) {
	var book Book
	row := db.QueryRow(
		"SELECT id, title, author, status FROM books WHERE id = $1", id)

	err := row.Scan(&book.Id, &book.Title, &book.Author, &book.Status)
	if err != nil {
		return Book{}, err
	}
	return book, nil
}

// updateBook — UPDATE a book's fields by id, return the updated book.
func updateBook(id int, book Book) (Book, error) {
	row := db.QueryRow(
		`UPDATE books SET title = $1, author = $2, status = $3
		 WHERE id = $4
		 RETURNING id, title, author, status`,
		book.Title, book.Author, book.Status, id)

	err := row.Scan(&book.Id, &book.Title, &book.Author, &book.Status)
	if err != nil {
		return Book{}, err
	}
	return book, nil
}

// deleteBook — DELETE a book by id.
// Uses db.Exec (not QueryRow) because nothing comes back.
// Checks RowsAffected to know if the book actually existed.
func deleteBook(id int) error {
	result, err := db.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("book with id %d not found", id)
	}
	return nil
}
