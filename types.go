package main

type Book struct {
	ID        int    `json:"id"`
	Titol     string `json:"titol"`
	Autor     string `json:"autor"`
	Prestatge string `json:"prestatge"`
	Posicio   string `json:"posicio"`
	Habitacio string `json:"habitacio"`
	Tipus     string `json:"tipus"`
	Editorial string `json:"editorial"`
	Idioma    string `json:"idioma"`
	Notes     string `json:"notes"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}