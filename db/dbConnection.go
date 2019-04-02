package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"strconv"
	"time"
)

// Exported object
var MgoSession *mgo.Session
var BuildingCollection string
var DatabaseName string
var UserLoginCollection string
var DemolishedCollection string
var logger *log.Logger

/**
Helper method to set DB properties
*/
func SetDbProperties(dbName string, bldngColl string, demolishedColl string, userColl string) {
	DatabaseName = dbName
	BuildingCollection = bldngColl
	DemolishedCollection = demolishedColl
	UserLoginCollection = userColl
}

func ConnectToDatabase(dbName string, dbHost string, dbUser string, dbPass string, dbTimeout string) {

	f, err := os.OpenFile("BuildingFootprintAnalyzer.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger = log.New(f, "prefix", log.LstdFlags)
	logger.Println("Starting mongodb session")

	timeout, err := strconv.Atoi(dbTimeout)

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{dbHost},
		Timeout:  time.Duration(timeout) * time.Second,
		Database: dbName,
		Username: dbUser,
		Password: dbPass,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		logger.Println("CreateSession: %s\n", err)
		panic(err)
	} else {
		MgoSession = mongoSession
		logger.Println("Created DB Session. Connected to:", dbName)
	}
}
