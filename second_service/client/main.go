package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"

	"serv2/handler"

	"github.com/go-chi/chi"
)

func main() {
	c, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		fmt.Println(err)
		return
	}

	handler := new(handler.Handler)
	handler.RpcClient = c

	r := chi.NewRouter()
	r.Post("/create-user", handler.CreateUser)
	r.Get("/get-user/{email}", handler.GetUser)

	if err := http.ListenAndServe("localhost:8080", r); err != nil {
		log.Fatal(err)
	}
}
