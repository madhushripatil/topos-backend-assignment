package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"topos-backend-assignment/db"
)

type Demolished struct {
	ID           bson.ObjectId `bson:"_id" json:"bin"`
	DoittID      int32         `bson:"doittID" json:"doittID"`
	Name         string        `bson:"name" json:"name"`
	Year         int16         `bson:"year" json:"year"`
	LastStatus   string        `bson:"lastStatus" json:"lastStatus"`
	FeatCode     int16         `bson:"featCode" json:"featCode"`
	HeightRoof   float32       `bson:"heightRoof" json:"heightRoof"`
	GroundLevel  int16         `bson:"groundLevel" json:"groundLevel"`
	ShapeArea    float32       `bson:"shapeArea" json:"shapeArea"`
	ShapeLength  float32       `bson:"shapeLength" json:"shapeLength"`
	LastModified time.Time     `bson:"lastModified" json:"lastModified"`
	GeomSource   string        `bson:"geomSource" json:"geomSource"`
	BaseBBL      int           `bson:"baseBBL" json:"baseBBL"`
	MPlutoBBL    int           `bson:"mplutoBBL" json:"mplutoBBL"`
}

func getDemolishedCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(db.DatabaseName).C(db.DemolishedCollection)
}

/**
Helper method to create a demolished structure archive
*/
func (demo *Demolished) CreateDemolishedStructure(session *mgo.Session, d Demolished) error {
	err := getDemolishedCollection(session).Insert(&d)
	return err
}

/**
Helper method to archive a demolished document
*/
func (demo *Demolished) ArchiveDemolishedBuildingFootPrint(session *mgo.Session, d Demolished) error {
	err := getDemolishedCollection(session).Insert(&d)
	return err
}
