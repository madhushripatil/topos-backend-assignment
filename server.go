package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Setting up server")
}

func main() {
	//defer db.MgoSession.Close()
	fmt.Println("Starting go service...")

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)

	}
}
