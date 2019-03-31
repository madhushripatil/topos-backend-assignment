package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"strconv"
	"time"
)

// Exported object
var MgoSession *mgo.Session
var hosts []string

func ConnectToDatabase(dbName string, dbHost string, dbUser string, dbPass string, dbTimeout string) {

	log.Println("Starting mongodb session")
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
		log.Fatalf("CreateSession: %s\n", err)
		panic(err)
	} else {
		MgoSession = mongoSession
		log.Println("Created DB Session. Connected to:", dbName)
	}
}
