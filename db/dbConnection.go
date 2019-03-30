package db

import (
	"gopkg.in/mgo.v2"
	"log"
)

// Exported objects
var MgoSession *mgo.Session

func ConnectToDatabase() {
	log.Println("Starting mongod session")
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	MgoSession = session
	//ensureIndexes()
}

func ensureIndexes() {
	ensureIndexOnBidder()
}

func ensureIndexOnBidder() {
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		Background: true,
	}

	err := MgoSession.DB("di_bidder_db").C("bidder").EnsureIndex(index)
	if err != nil {
		log.Fatal("Unique index does not exist on collection : bidder")
		panic("Unique index does not exist on collection : bidder")
	}
}
