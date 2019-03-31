package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION = "building"
const DBNAME = "buildingFootprintsDB"

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

func getDBCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(DBNAME).C(COLLECTION)
}

func findById(session *mgo.Session, id string) (Building, error) {
	var bld Building
	err := getDBCollection(session).FindId(bson.ObjectIdHex(id)).One(&bld)
	return bld, err
}

func (building *Building) GetAllBuildingFootPrints(session *mgo.Session) ([]Building, error) {
	var buildings []Building
	err := getDBCollection(session).Find(bson.M{}).All(&buildings)
	return buildings, err
}

func (building *Building) GetAllBuildingsCount(session *mgo.Session) ([]Building, error, int) {
	var buildings []Building
	var i interface{}
	cnt, err := getDBCollection(session).Find(i).Count()
	return buildings, err, cnt
}

func (building *Building) FindBuildingFootPrintById(session *mgo.Session, id string) (Building, error) {
	var bld Building
	bld, err := findById(session, id)
	return bld, err
}

func (building *Building) CreateBuildingFootPrint(session *mgo.Session, b Building) error {
	err := getDBCollection(session).Insert(&b)
	return err
}

func (building *Building) DeleteBuildingFootPrint(session *mgo.Session, b Building) error {
	err := getDBCollection(session).Remove(&b)
	return err
}

func (building *Building) UpdateBuildingFootPrint(session *mgo.Session, b Building) error {
	err := getDBCollection(session).UpdateId(b.ID, &b)
	return err
}
