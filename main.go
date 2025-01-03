package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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

func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.Data)

	data := string(b.Pos) + b.Timestamp + string(bytes) + b.PrevHash

	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

func CreateBlock(prevBlock *Block, data BookCheckout) *Block {
	block := &Block{
		Pos:       prevBlock.Pos + 1,
		Timestamp: time.Now().String(),
		PrevHash:  prevBlock.Hash,
	}

	block.generateHash()
	return block
}

func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	if b.Hash != hash {
		return false
	}
	return true
}

func validBlock(newBlock, prevBlock *Block) bool {
	if prevBlock.Hash != newBlock.PrevHash {
		return false
	}

	if !newBlock.validateHash(newBlock.Hash) {
		return false
	}

	if prevBlock.Pos+1 != newBlock.Pos {
		return false
	}

	return true
}

func (bc *Blockchain) AddBlock(data BookCheckout) {

	prevBlock := bc.block[len(bc.block)-1]

	block := CreateBlock(prevBlock, data)

	if validBlock(block, prevBlock) {
		bc.block = append(bc.block, block)
	}

}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(BlockChain.block, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	io.WriteString(w, string(jbytes))
}

func newBook(w http.ResponseWriter, r *http.Request) {

	var book Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %s", err)
		w.Write([]byte("Error decoding book"))
		return
	}

	h := md5.New()
	io.WriteString(h, book.Isbn+book.PublishDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	// Convert to JSON
	resp, err := json.MarshalIndent(book, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func writeBlock(w http.ResponseWriter, r *http.Request) {

	var bookCheckout BookCheckout
	if err := json.NewDecoder(r.Body).Decode(&bookCheckout); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %s", err)
		w.Write([]byte("Could no write block"))
		return
	}

	BlockChain.AddBlock(bookCheckout)

	// Convert to JSON
	resp, err := json.MarshalIndent(bookCheckout, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

var BlockChain *Blockchain

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

func NewBlockchain() *Blockchain {

	// Slice of mukltiple blocks
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func main() {

	BlockChain := NewBlockchain()

	r := mux.NewRouter()

	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	go func() {

		for _, block := range BlockChain.block {
			fmt.Printf("Prev. hash: %x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data: %v\n", string(bytes))
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Println()
		}

	}()

	log.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r).Error())
}
