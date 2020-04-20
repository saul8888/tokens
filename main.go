package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tokens/authentication"
	//"github.com/tokens/authentication"
)

func main() {

	r := chi.NewRouter()

	r.Post("/validate", authentication.ValidateToken)
	r.Post("/login", authentication.Login)
	//r.Get("/login", prueba)
	//r.Post("/validate", prueba)

	log.Println("listen on port 8000")
	http.ListenAndServe(":8000", r)
}

//Bearer {token}

//func prueba(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("hi"))
//}
