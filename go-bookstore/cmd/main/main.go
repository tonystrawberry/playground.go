package main

import (
	"log" // log error
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tonystrawberry/go-bookstore/pkg/routes"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:8000", r))
}