package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"strconv"
	"topos-backend-assignment/db"
	"topos-backend-assignment/models"
)

type ResponseMessage struct {
	Status   int    `json:"status"`
	ErrorMsg string `json:"error"`
	Message  string `json:"message"`
}

var DbName string
var CollectionName string
var ServerPort string
var DBHost string
var DBUsername string
var DBPassword string
var DBTimeout string

/**
API URL - http://localhost:8000/
Method	- GET
Params	- None
This API returns the count of all the BuildingFootPrint Data in the Database
*/
func AllBuildingsCount(writer http.ResponseWriter, req *http.Request) {
	var bld models.Building
	var js []byte
	var err error
	var msg ResponseMessage
	var cnt int

	err, cnt = bld.GetAllBuildingsCount(db.MgoSession)

	if err != nil {
		msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error getting Building count"}
		js, err = json.Marshal(msg)

	} else {
		msg = ResponseMessage{Status: http.StatusOK, ErrorMsg: "", Message: "Total Building Footprints " + strconv.Itoa(cnt)}
		js, err = json.Marshal(msg)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://localhost:8000/buildingFootprints
Method	- GET
Params	- None
This API returns all the BuildingFootPrint Data from the Database
*/
func AllBuildingFootPrints(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	var js []byte
	var err error
	var buildings []models.Building
	var msg ResponseMessage

	buildings, err = bld.GetAllBuildingFootPrints(db.MgoSession)

	if err != nil {
		msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error getting FootPrint Data"}
		js, err = json.Marshal(msg)

	} else {
		if len(buildings) != 0 {
			js, err = json.Marshal(buildings)
		} else {
			msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "Empty Dataset"}
			js, err = json.Marshal(msg)
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://localhost:8000/buildingFootprints/{id}
Method	- GET
Params	- id : Building footprint ID
This API returns the BuildingFootPrint Data for a specific ID
*/
func GetBuildingFootPrintsById(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	params := mux.Vars(request)
	var js []byte
	var err error
	var msg ResponseMessage

	bld, err = bld.FindBuildingFootPrintById(db.MgoSession, params["id"])
	if err != nil {
		msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error getting FootPrint Data"}
		js, err = json.Marshal(msg)
	} else {
		if (models.Building{}) != bld {
			js, err = json.Marshal(bld)
		} else {
			msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "Empty Dataset"}
			js, err = json.Marshal(msg)
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://localhost:8000/buildingFootprints/{id}
Method	- DELETE
Params	- id : Building footprint ID
This API deletes the BuildingFootPrint Data for a specific ID
*/
func DeleteBuildingFootPrints(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	params := mux.Vars(request)
	var js []byte
	var err error
	var msg ResponseMessage

	bld, err = bld.FindBuildingFootPrintById(db.MgoSession, params["id"])
	if err != nil {
		msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Deleting FootPrint Data"}
		js, err = json.Marshal(msg)
	} else {
		if (models.Building{}) != bld {
			if err = bld.DeleteBuildingFootPrint(db.MgoSession, bld); err != nil {
				msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Deleting FootPrint Data"}
				js, err = json.Marshal(msg)
			} else {
				msg = ResponseMessage{Status: http.StatusOK, ErrorMsg: "None", Message: "BuildingFootPrint successfully deleted!"}
				js, err = json.Marshal(msg)
			}

		} else {
			msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "No Such BuildingFootPrint"}
			js, err = json.Marshal(msg)
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://localhost:8000/buildingFootprints/{id}
Method	- PUT
Params	- id : Building footprint ID
Request Body - JSON Object containing all the fields to be updated
This API updates the BuildingFootPrint Data for a specific ID
*/
func UpdateBuildingFootPrints(writer http.ResponseWriter, request *http.Request) {

	var bld models.Building
	params := mux.Vars(request)
	var js []byte
	var err error
	var msg ResponseMessage

	bld, err = bld.FindBuildingFootPrintById(db.MgoSession, params["id"])
	if err != nil {
		msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Updating FootPrint Data"}
		js, err = json.Marshal(msg)
	} else {
		if (models.Building{}) != bld {
			if err = bld.UpdateBuildingFootPrint(db.MgoSession, bld); err != nil {
				msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Updating FootPrint Data"}
				js, err = json.Marshal(msg)
			} else {
				msg = ResponseMessage{Status: http.StatusOK, ErrorMsg: "None", Message: "BuildingFootPrint successfully updated!"}
				js, err = json.Marshal(msg)
			}

		} else {
			msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "No Such BuildingFootPrint"}
			js, err = json.Marshal(msg)
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://localhost:8000/buildingFootprints
Method	- POST
Params	- None
Request Body - JSON Object containing all the fields
This API deletes the BuildingFootPrint Data for a specific ID
*/
func AddBuildingFootPrints(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	var js []byte
	var err error
	var msg ResponseMessage
	defer request.Body.Close()

	if err = json.NewDecoder(request.Body).Decode(&bld); err != nil {
		msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Creating FootPrint Data"}
		js, err = json.Marshal(msg)
	} else {
		bld.ID = bson.NewObjectId()
		if err = bld.CreateBuildingFootPrint(db.MgoSession, bld); err != nil {
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Creating FootPrint Data"}
			js, err = json.Marshal(msg)
		} else {
			msg = ResponseMessage{Status: http.StatusCreated, ErrorMsg: "None", Message: "BuildingFootPrint successfully created!"}
			js, err = json.Marshal(msg)
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://localhost:8000/buildingFootprints/type/{bldType}
Method	- GET
Params	- None
This API returns all the Buildings by Type
*/
func GetBuildingsByType(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	params := mux.Vars(request)
	var js []byte
	var err error
	var buildings []models.Building
	var msg ResponseMessage
	var bldngType int

	bldngType, err = strconv.Atoi(params["bldType"])
	buildings, err = bld.FindAllBuildingsByType(db.MgoSession, bldngType)
	if err != nil {
		msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error getting Buildings by type"}
		js, err = json.Marshal(msg)
	} else {
		if len(buildings) != 0 {
			js, err = json.Marshal(buildings)
		} else {
			msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "Empty Dataset"}
			js, err = json.Marshal(msg)
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
Main method - Execution starts here
*/
func main() {
	defer db.MgoSession.Close()
	fmt.Println("Starting BuildingFootprint services...")

	// Load env file
	e := godotenv.Load("buildingFootprint.env")
	if e != nil {
		fmt.Print(e)
	}

	DbName = os.Getenv("db_name")
	CollectionName = os.Getenv("collection_name")
	DBHost = os.Getenv("db_host")
	ServerPort = os.Getenv("server_port")
	DBTimeout = os.Getenv("db_timeout")
	DBUsername = os.Getenv("db_username")
	DBPassword = os.Getenv("db_pass")

	port := fmt.Sprintf(":%s", ServerPort)

	route := mux.NewRouter()
	route.HandleFunc("/", AllBuildingsCount).Methods("GET")
	route.HandleFunc("/buildingFootprints", AddBuildingFootPrints).Methods("POST")
	route.HandleFunc("/buildingFootprints/{id}", UpdateBuildingFootPrints).Methods("PUT")
	route.HandleFunc("/buildingFootprints/{id}", DeleteBuildingFootPrints).Methods("DELETE")
	route.HandleFunc("/buildingFootprints", AllBuildingFootPrints).Methods("GET")
	route.HandleFunc("/buildingFootprints/{id}", GetBuildingFootPrintsById).Methods("GET")
	route.HandleFunc("/buildingFootprints/type/{bldType}", GetBuildingsByType).Methods("GET")

	// The following function call makes a database connection
	db.ConnectToDatabase(DbName, DBHost, DBUsername, DBPassword, DBTimeout)
	models.SetDbProperties(CollectionName, DbName)

	// The Server listens on port for incoming requests and routes requests
	if err := http.ListenAndServe(port, route); err != nil {
		log.Fatal(err)
	}
}
