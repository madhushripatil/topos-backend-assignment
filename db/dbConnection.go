package db

import (
	"gopkg.in/mgo.v2"
	"log"
)

// Exported object
var MgoSession *mgo.Session

func ConnectToDatabase(serverHost string) {
	log.Println("Starting mongodb session")
	session, err := mgo.Dial(serverHost)
	if err != nil {
		panic(err)
	}
	MgoSession = session
}
