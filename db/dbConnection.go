package db

import (
	"gopkg.in/mgo.v2"
	"log"
)

// Exported object
var MgoSession *mgo.Session

func ConnectToDatabase() {
	log.Println("Starting mongodb session")
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	MgoSession = session
}
