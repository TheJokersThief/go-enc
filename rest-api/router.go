package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// our main function
func main() {
	// https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo
	router := mux.NewRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}
