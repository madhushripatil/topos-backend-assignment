package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"topos-backend-assignment/db"
)

func AllBuildingsCount(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(writer, "Setting up server")
}

func AddBuildingFootPrints(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(writer, "Setting up server")
}

func UpdateBuildingFootPrints(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(writer, "Setting up server")
}

func DeleteBuildingFootPrints(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(writer, "Setting up server")
}

func AllBuildingFootPrints(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(writer, "Setting up server")
}

func main() {
	defer db.MgoSession.Close()
	fmt.Println("Starting go service...")

	route := mux.NewRouter()
	route.HandleFunc("/", AllBuildingsCount).Methods("GET")
	route.HandleFunc("/buildingFootprints", AddBuildingFootPrints).Methods("POST")
	route.HandleFunc("/buildingFootprints/{id}", UpdateBuildingFootPrints).Methods("PUT")
	route.HandleFunc("/buildingFootprints/{id}", DeleteBuildingFootPrints).Methods("DELETE")
	route.HandleFunc("/buildingFootprints", AllBuildingFootPrints).Methods("GET")

	db.ConnectToDatabase()

	if err := http.ListenAndServe(":8000", route); err != nil {
		log.Fatal(err)
	}
}
