package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"topos-backend-assignment/db"
)

type Borough struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	BoroughCode int64         `bson:"boroughCode" json:"boroughCode"`
	BoroughName string        `bson:"boroughName" json:"boroughName"`
}

func getBoroughCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(db.DatabaseName).C(db.BoroughTypeCollection)
}

/**
Helper method to find Building Feat Code by type
*/
func GetBoroughCodeByName(session *mgo.Session, borough string) (int, error) {
	var b Borough
	err := getBoroughCollection(session).Find(bson.M{"boroughName": borough}).One(&b)
	return int(b.BoroughCode), err
}
