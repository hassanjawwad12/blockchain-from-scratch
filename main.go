package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Block struct {
}

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publish_date"`
	Isbn        string `json:"isbn"`
}

type Blockchain struct {
	// Slice of multiple blocks
	block []*Block
}

type BookCheck struct {
}

func getBlock(w http.ResponseWriter, r *http.Request) {
}

var Blockchain *Blockchain

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", getBlock).Methods("GET")
	// r.HandleFunc("/", writeBlock).Methods("POST")
	// r.HandleFunc("/new", newBook).Methods("POST")

	log.Println("Listening on port 3000")
	log.Fatalf(http.ListenAndServe(":3000", r).Error())
}
