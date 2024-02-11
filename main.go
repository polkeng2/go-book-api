package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS books (id SERIAL PRIMARY KEY, title TEXT, author TEXT)")

	fmt.Println("Server is running on port 8000")

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks(db)).Methods("GET")
	router.HandleFunc("/books", createBook(db)).Methods("POST")
	router.HandleFunc("/books", updateBook(db)).Methods("PUT")
	router.HandleFunc("/books", deleteBook(db)).Methods("DELETE")

	http.ListenAndServe(":8000", handleMiddleware(router))
}


func handleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM books")
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		defer rows.Close()

		books := make([]Book, 0)
		for rows.Next() {
			book := Book{}
			err := rows.Scan(&book.ID, &book.Title, &book.Author)
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}
			books = append(books, book)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(books)
	}
}

func createBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		book := Book{}
		err := json.NewDecoder(r.Body).Decode(&book)
		fmt.Println(book)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}


		_, err = db.Exec("INSERT INTO books (title, author) VALUES ($1, $2)", book.Title, book.Author)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func updateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		book := Book{}
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		_, err = db.Exec("UPDATE books SET title = $1, author = $2 WHERE id = $3", book.Title, book.Author, book.ID)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func deleteBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		book := Book{}
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		_, err = db.Exec("DELETE FROM books WHERE id = $1", book.ID)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}