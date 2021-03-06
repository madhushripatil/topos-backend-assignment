package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"topos-backend-assignment/controller"
	"topos-backend-assignment/db"
)

// DbName exported
var DbName string

// BuildingCollectionName exported
var BuildingCollectionName string

// DemolishedCollectionName exported
var DemolishedCollectionName string

// ServerPort exported
var ServerPort string

// DBHost exported
var DBHost string

// DBUsername exported
var DBUsername string

// DBPassword exported
var DBPassword string

// DBTimeout exported
var DBTimeout string

// UserLoginCollectionName exported
var UserLoginCollectionName string

// BldngFeatTypeCollectionName exported
var BldngFeatTypeCollectionName string

// BoroughCollectionName exported
var BoroughCollectionName string

// Logger exported
var Logger *log.Logger

// JwtKey exported
var JwtKey []byte

/**
Main method - Execution starts here
*/
func main() {

	F, err := os.OpenFile("BuildingFootprintAnalyzer.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer F.Close()

	Logger = log.New(F, "prefix", log.LstdFlags)

	defer db.MgoSession.Close()

	// Load env file
	Logger.Println("Loading the configuration file...")
	envFile := fmt.Sprintf("%s.env", os.Getenv("BUILDING_ENV"))
	e := godotenv.Load(envFile)

	if e != nil {
		// close the application if the configuration file open fails
		Logger.Println("Error loading the config file, Closing the application.")
		panic(e)
	} else {

		Logger.Println("Starting BuildingFootprint services...")

		// Connect to the database and start services
		DbName = os.Getenv("db_name")
		BuildingCollectionName = os.Getenv("bldng_collection_name")
		DemolishedCollectionName = os.Getenv("demolished_collection_name")
		BldngFeatTypeCollectionName = os.Getenv("bldng_feat_type_collection_name")
		BoroughCollectionName = os.Getenv("borough_collection_name")
		DBHost = os.Getenv("db_host")
		ServerPort = os.Getenv("server_port")
		DBTimeout = os.Getenv("db_timeout")
		DBUsername = os.Getenv("db_username")
		DBPassword = os.Getenv("db_pass")
		UserLoginCollectionName = os.Getenv("userlogin_collection_name")
		JwtKey = []byte(os.Getenv("jwt_key"))

		port := fmt.Sprintf(":%s", ServerPort)

		route := mux.NewRouter()
		route.HandleFunc("/", controller.AllBuildingsCount).Methods("GET")
		route.HandleFunc("/buildingFootprints", controller.AddBuildingFootPrints).Methods("POST")
		route.HandleFunc("/buildingFootprints/{id}", controller.UpdateBuildingFootPrints).Methods("PUT")
		route.HandleFunc("/buildingFootprints/{id}", controller.DeleteBuildingFootPrints).Methods("DELETE")
		route.HandleFunc("/buildingFootprints", controller.AllBuildingFootPrints).Methods("GET")
		route.HandleFunc("/buildingFootprints/{id}", controller.GetBuildingFootPrintsById).Methods("GET")
		route.HandleFunc("/buildingFootprints/type/{bldType}", controller.GetBuildingsByType).Methods("GET")
		route.HandleFunc("/buildingFootprints/borough/{boroughName}", controller.GetBuildingsByBorough).Methods("GET")
		route.HandleFunc("/buildingFootprints/buildingHeightAndArea/{minHeight}/{minArea}", controller.GetTallAndWideBuildings).Methods("GET")
		route.HandleFunc("/buildingFootprints/demolishedStructuresByConstructedYear/{year}", controller.GetAllDemolishedStructuresByYear).Methods("GET")

		route.HandleFunc("/authenticate/login", controller.LoginUser).Methods("POST")
		route.HandleFunc("/authenticate/signup", controller.SignUp).Methods("POST")

		route.HandleFunc("/transform/archiveDemolishedStructures", controller.ArchiveAllDemolishedStructures).Methods("POST")

		// The following function call makes a database connection
		db.ConnectToDatabase(DbName, DBHost, DBUsername, DBPassword, DBTimeout)
		db.SetDbProperties(DbName, BuildingCollectionName, DemolishedCollectionName, UserLoginCollectionName, BldngFeatTypeCollectionName, BoroughCollectionName)
		controller.UseLogger(Logger)
		controller.SetJWTSecret(JwtKey)

		// The Server listens on port for incoming requests and routes requests
		if err := http.ListenAndServe(port, route); err != nil {
			Logger.Println("Error listening on port ", port)
			panic(err)
		}
		Logger.Println("Server Listening on port... ", port)
	}
}
