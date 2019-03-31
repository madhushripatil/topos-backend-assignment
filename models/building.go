package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var CollectionName string
var DbName string

type Buildings struct {
	Buildings []Building `json:"buildings"`
}

type Building struct {
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

/**
Helper method to set DB properties
*/
func SetDbProperties(c string, d string) {
	CollectionName = c
	DbName = d
}

func getDBCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(DbName).C(CollectionName)
}

/**
Helper method to find a document by ID
*/
func findById(session *mgo.Session, id string) (Building, error) {
	var bld Building
	err := getDBCollection(session).FindId(bson.ObjectIdHex(id)).One(&bld)
	return bld, err
}

/**
Helper method to find all documents in a collection
*/
func (building *Building) GetAllBuildingFootPrints(session *mgo.Session) ([]Building, error) {
	var buildings []Building
	err := getDBCollection(session).Find(bson.M{}).All(&buildings)
	return buildings, err
}

/**
Helper method to find count of all documents in a collection
*/
func (building *Building) GetAllBuildingsCount(session *mgo.Session) (error, int) {
	var i interface{}
	cnt, err := getDBCollection(session).Find(i).Count()
	return err, cnt
}

/**
Helper method to find a document by ID
*/
func (building *Building) FindBuildingFootPrintById(session *mgo.Session, id string) (Building, error) {
	var bld Building
	if bson.IsObjectIdHex(id) {
		bld, err := findById(session, id)
		return bld, err
	} else {
		return bld, errors.New("please provide a valid Object ID")
	}
}

/**
Helper method to create a new document
*/
func (building *Building) CreateBuildingFootPrint(session *mgo.Session, b Building) error {
	err := getDBCollection(session).Insert(&b)
	return err
}

/**
Helper method to delete a specific document
*/
func (building *Building) DeleteBuildingFootPrint(session *mgo.Session, b Building) error {
	err := getDBCollection(session).Remove(&b)
	return err
}

/**
Helper method to update a specific document
*/
func (building *Building) UpdateBuildingFootPrint(session *mgo.Session, b Building) error {
	err := getDBCollection(session).UpdateId(b.ID, &b)
	return err
}

/**
Helper method to find all Buildings by type
*/
func (building *Building) FindAllBuildingsByType(session *mgo.Session, bldType int) ([]Building, error) {
	var buildings []Building
	err := getDBCollection(session).Find(bson.M{"featCode": bldType}).All(&buildings)
	return buildings, err
}
