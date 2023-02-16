package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type Reply struct {
	Salt string `json:"salt"`
}

func Generate(w http.ResponseWriter, r *http.Request) {
	var resp Reply
	resp.Salt = randSeq(12)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/generate-salt", Generate)

	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}
