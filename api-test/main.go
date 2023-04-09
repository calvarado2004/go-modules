package main

import (
	"log"
	"net/http"
)

func main() {

	mux := routes()

	log.Println("Listening on port 8080...")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", login)
	mux.HandleFunc("/api/logout", logout)

	return mux
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login"))
}

func logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("logout"))
}
