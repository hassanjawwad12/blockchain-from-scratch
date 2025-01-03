package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Block struct {
	Pos       int
	Data      BookCheckout
	Timestamp string
	Hash      string
	PrevHash  string
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

type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
}

var blockchain *Blockchain

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", getBlockchain).Methods("GET")
	// r.HandleFunc("/", writeBlock).Methods("POST")
	// r.HandleFunc("/new", newBook).Methods("POST")

	log.Println("Listening on port 3000")
	log.Fatalf(http.ListenAndServe(":3000", r).Error())
}
