package main

import (
	"log"
	"net/http"
	"toolkit"
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
	mux.HandleFunc("/simulated-service", simulatedService)

	// print a message to the console
	println("Listening on port 8081")

	// start the server

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}

// receivePost is a handler function that receives a POST request
func receivePost(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	var t toolkit.Tools

	// decode the request body into the requestPayload struct
	err := t.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = t.ErrorJSON(w, err)
		return
	}

	responsePayload := ResponsePayload{
		Message: "you hit the handler function",
	}

	err = t.WriteJSON(w, http.StatusAccepted, responsePayload)
	if err != nil {
		_ = t.ErrorJSON(w, err)
		return
	}

}

// remoteService is a handler function that receives a POST request
func remoteService(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	var t toolkit.Tools

	// decode the request body into the requestPayload struct
	err := t.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = t.ErrorJSON(w, err)
		return
	}

	// call the simulated service
	_, statusCode, err := t.PushJSONToRemote("http://localhost:8081/simulated-service", requestPayload)
	if err != nil {
		_ = t.ErrorJSON(w, err)
		return
	}

	responsePayload := ResponsePayload{
		Message:    "you hit the handler function",
		StatusCode: statusCode,
	}

	err = t.WriteJSON(w, http.StatusAccepted, responsePayload)
	if err != nil {
		_ = t.ErrorJSON(w, err)
		return
	}

}

// simulatedService is a handler function that receives a POST request
func simulatedService(w http.ResponseWriter, r *http.Request) {

	var payload ResponsePayload

	payload.Message = "simulated service"
	payload.StatusCode = http.StatusAccepted

	var t toolkit.Tools

	err := t.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		_ = t.ErrorJSON(w, err)
		return
	}

}
