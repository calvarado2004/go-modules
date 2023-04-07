package main

import (
	"log"
	"net/http"
)

type RequestPayload struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type ResponsePayload struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code,omitempty"`
}

func main() {

	// create a new serve mux
	mux := http.NewServeMux()

	// register the handler functions
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))

	mux.HandleFunc("/receive-post", receivePost)
	mux.HandleFunc("/remote-service", remoteService)

	// print a message to the console
	println("Listening on port 8081")

	// start the server

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}

func receivePost(w http.ResponseWriter, r *http.Request) {

}

func remoteService(w http.ResponseWriter, r *http.Request) {

}
