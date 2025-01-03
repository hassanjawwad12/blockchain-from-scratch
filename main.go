package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", getBlock).Methods("GET")
	// r.HandleFunc("/", writeBlock).Methods("POST")
	// r.HandleFunc("/new", newBook).Methods("POST")

	log.Println("Listening on port 3000")
	log.Fatalf(http.ListenAndServe(":3000", r).Error())
}
