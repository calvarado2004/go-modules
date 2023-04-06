package main

import (
	"log"
	"net/http"
	"toolkit"
)

// main is the entry point for the application
func main() {
	// get some routes

	mux := routes()

	// start a server
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}

// routes returns a http.Handler that handles all the routes
func routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))

	mux.HandleFunc("/download", downloadFile)

	return mux
}

// downloadFile downloads the file
func downloadFile(w http.ResponseWriter, r *http.Request) {

	t := toolkit.Tools{}
	t.DownloadStaticFile(w, r, "./files", "timesquare.jpeg", "timesquare.jpeg")
}
