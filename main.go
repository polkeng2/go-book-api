package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()	

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS books (id SERIAL PRIMARY KEY, titol TEXT, autor TEXT, prestatge TEXT, posicio TEXT, habitacio TEXT, tipus TEXT, editorial TEXT, idioma TEXT, notes TEXT)")

	if err != nil {
		log.Fatal(err)
	}
	
	importAndParseDB(db)

	router := mux.NewRouter()

	router.HandleFunc("/", returnHelloWorld()).Methods("GET")
	router.HandleFunc("/books/first", getFirstBook(db)).Methods("GET")
 	router.HandleFunc("/books", getBooks(db)).Methods("GET")
	router.HandleFunc("/books", createBook(db)).Methods("POST")
	router.HandleFunc("/books/{id}", deleteBook(db)).Methods("DELETE") 
	router.HandleFunc("/books/{id}", updateBook(db)).Methods("PUT") 
	router.HandleFunc("/login", handleLogin()).Methods("POST")

	enhancedRouter := enableCORS(handleMiddleware(router))

	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), enhancedRouter))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Check if the request is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})

}


func handleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if(r.Method != "GET" && r.URL.Path != "/login") {
			token := r.Header.Get("Authorization")
			if token != os.Getenv("TOKEN") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userData := LoginData{}
		err := json.NewDecoder(r.Body).Decode(&userData)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		
		if userData.Email == os.Getenv("EMAIL") && userData.Password == os.Getenv("PASSWORD") {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(os.Getenv("TOKEN"))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func returnHelloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello World")
	}
}

func getFirstBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		book := Book{}
		err := db.QueryRow("SELECT * FROM books LIMIT 1").Scan(&book.ID, &book.Titol, &book.Autor, &book.Prestatge, &book.Posicio, &book.Habitacio, &book.Tipus, &book.Editorial, &book.Idioma, &book.Notes)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(book)
	}
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
			err := rows.Scan(&book.ID, &book.Titol, &book.Autor, &book.Prestatge, &book.Posicio, &book.Habitacio, &book.Tipus, &book.Editorial, &book.Idioma, &book.Notes)
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
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		book := Book{}
		err := json.NewDecoder(r.Body).Decode(&book)
		fmt.Println(book)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		_, err = db.Exec("INSERT INTO books (titol, autor, prestatge, posicio, habitacio, tipus, editorial, idioma, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", book.Titol, book.Autor, book.Prestatge, book.Posicio, book.Habitacio, book.Tipus, book.Editorial, book.Idioma, book.Notes)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
	}
}

func updateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		book := Book{}
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		//TODO: check id always starts at 1
		_, err = db.Exec("UPDATE books SET titol = $1, autor = $2, prestatge = $3, posicio = $4, habitacio = $5, tipus = $6, editorial = $7, idioma = $8, notes = $9 WHERE id = $10", book.Titol, book.Autor, book.Prestatge, book.Posicio, book.Habitacio, book.Tipus, book.Editorial, book.Idioma, book.Notes, id)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func deleteBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {		
		id := mux.Vars(r)["id"]
		_, err := db.Exec("DELETE FROM books WHERE id = $1", id)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		_, err = db.Exec("SELECT setval('books_id_seq', (SELECT MAX(id) FROM books))")
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func importAndParseDB(db *sql.DB) {

	// Check the db is empty
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		return
	}

	// Open the file
	csvfile, err := os.Open("biblioteca.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		
		//make sure the tipus section starts with capital letter
		record[5] = strings.Title(record[5])
		record[6] = strings.Title(record[6])
		if record[8] != "" {
			record[8] = strings.Title(record[8])
		}
		
		_, dbErr := db.Exec("INSERT INTO books (titol, autor, prestatge, posicio, habitacio, tipus, editorial, idioma, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7], record[8])
		if dbErr != nil {
			log.Fatal(dbErr)
		}
	}
}