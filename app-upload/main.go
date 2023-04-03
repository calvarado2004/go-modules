package main

import (
	"fmt"
	"log"
	"net/http"
	"toolkit"
)

// main function to boot up everything
func main() {

	mux := routes()

	log.Println("Server is running on port 8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}

}

// routes function to handle all routes
func routes() http.Handler {

	mux := http.NewServeMux()

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))

	mux.HandleFunc("/upload", uploadFiles)

	mux.HandleFunc("/upload-single/", uploadFile)

	return mux
}

// uploadFiles function to handle multiple file uploads
func uploadFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	t := toolkit.Tools{
		MaxFileSize:      1024 * 1024 * 1024,
		AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif"},
	}

	files, err := t.UploadFiles(r, "./uploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := ""
	for _, file := range files {
		out += fmt.Sprintf("Uploaded to: %s, renamed to: %s\n", file.OriginalFile, file.NewFileName)
	}

	_, _ = w.Write([]byte(out))

}

// uploadFile function to handle single file uploads
func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	t := toolkit.Tools{
		MaxFileSize:      1024 * 1024 * 1024,
		AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif"},
	}

	file, err := t.UploadOneFile(r, "./uploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := fmt.Sprintf("Uploaded to: %s, renamed to: %s\n", file.OriginalFile, file.NewFileName)

	_, _ = w.Write([]byte(out))
}
