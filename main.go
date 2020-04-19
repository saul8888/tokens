package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tokens/authentication"
)

func main() {

	r := chi.NewRouter()

	r.Post("/login", authentication.login)
	r.Post("/validate", authentication.ValidateToken)

	log.Println("listen on port 8000")
	http.ListenAndServe(":8000", r)
}
