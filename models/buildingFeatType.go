package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"topos-backend-assignment/db"
)

type BuildingFeatType struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	FeatCode int16         `bson:"featCode" json:"featCode"`
	FeatType string        `bson:"featType" json:"featType"`
}

func getBldngFeatTypeCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(db.DatabaseName).C(db.BuildingFeatTypeCollection)
}

/**
Helper method to find Building Feat Code by type
*/
func GetFeatCodeByName(session *mgo.Session, bldType string) (int, error) {
	var b BuildingFeatType
	err := getBldngFeatTypeCollection(session).Find(bson.M{"featType": bldType}).One(&b)
	return int(b.FeatCode), err
}
