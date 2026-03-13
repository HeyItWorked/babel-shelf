// main.go — Entry point. Connects to Postgres, sets up routes, starts the server.
package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Postgres driver — underscore means "import for side effects only"
	// (registers itself with database/sql, we never call it directly)
)

// db is the database connection, shared across handlers.
// package-level variable — all files in package main can access it (db.go, handlers.go, etc).
// lowercase = not exported to other packages, but visible within this package.
var db *sql.DB

// setupRouter creates the mapping table — which URL goes to which handler
func setupRouter() http.Handler {
	// set up routes here
	// POST /books    → CreateBook
	// GET /books     → GetBooks
	// GET /books/{id} → GetBook
	// PUT /books/{id} → UpdateBook
	// DELETE /books/{id} → DeleteBook

	// router/request mapper
	mux := http.NewServeMux()

	// HandleFunc maps the handler func to REST call
	mux.HandleFunc("POST /books", CreateBook)
	mux.HandleFunc("GET /books", GetBooks)
	mux.HandleFunc("GET /books/{id}", GetBook)
	mux.HandleFunc("PUT /books/{id}", UpdateBook)
	mux.HandleFunc("DELETE /books/{id}", DeleteBook)

	return mux
}

func main() {
	// read DATABASE_URL from environment variable
	var err error
	databaseUrl := os.Getenv("DATABASE_URL")

	if databaseUrl == "" {
		fmt.Println("DATABASE_URL is required")
  		os.Exit(1)
	}

	// connect to Postgres using sql.Open
	// When you imported _"github.com/lib/pq", that package registered itself 
	// under the name "postgres"
	db, err = sql.Open("postgres", databaseUrl)
	if err != nil {
		fmt.Printf("connection failed: %v\n", err)
  		os.Exit(1)
	}
	defer db.Close()

	// verify connection works using db.Ping
	// pings like how we check 8.8.8.8 but for db
	err = db.Ping()
	if err != nil {
		fmt.Printf("could not ping database: %v\n", err)
  		os.Exit(1)	
	}
	fmt.Println("connected to database")

	// start the server on :8080
	router := setupRouter()
	fmt.Println("server running on port 8080")

	// ListenAndServe blocks forever while the server runs.
	// It only returns when something goes wrong.
	err = http.ListenAndServe(":8080", router)
	fmt.Println("server failed:", err)
	os.Exit(1)
}
