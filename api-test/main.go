package main

import (
	"log"
	"net/http"
	"toolkit"
)

func main() {

	mux := routes()

	log.Println("Listening on port 8080...")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

// routes returns a new ServeMux
func routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	mux.HandleFunc("/api/login", login)
	mux.HandleFunc("/api/logout", logout)

	return mux
}

// login and logout handlers
func login(w http.ResponseWriter, r *http.Request) {

	var tools toolkit.Tools

	var payload struct {
		Email    string `json:"username"`
		Password string `json:"password"`
	}

	err := tools.ReadJSON(w, r, &payload)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	var resPayload toolkit.JSONResponse

	if payload.Email == "me@here.com" && payload.Password == "password" {
		resPayload = toolkit.JSONResponse{
			Msg: "Logged in",
		}
		_ = tools.WriteJSON(w, http.StatusAccepted, resPayload)
		return
	} else {
		resPayload = toolkit.JSONResponse{
			Msg:   "Invalid login",
			Error: true,
		}
		_ = tools.WriteJSON(w, http.StatusUnauthorized, resPayload)
		return
	}

}

func logout(w http.ResponseWriter, r *http.Request) {

	var tools toolkit.Tools

	payload := toolkit.JSONResponse{
		Msg: "Logged out",
	}

	_ = tools.WriteJSON(w, http.StatusAccepted, payload)
}
