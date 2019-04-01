package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
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
func (building *Building) UpdateBuildingFootPrint(session *mgo.Session, b Building, updBldng Building) error {

	v := reflect.ValueOf(updBldng)
	filter := bson.M{"_id": building.ID}
	updateMap := map[string]interface{}{}

	for i := 0; i < v.NumField(); i++ {
		flag := isEmpty(v.Field(i).Kind(), v.Field(i))
		nm := v.Type().Field(i).Tag.Get("json")

		if !flag {
			// add non empty values to update
			updateMap[nm] = v.Field(i).Interface()
		}
	}

	// create a map of values to update
	incDoc := bson.M{}
	for k, v := range updateMap {
		incDoc[k] = v
	}
	change := bson.M{"$set": incDoc}

	err := getDBCollection(session).Update(filter, change)
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

/**
Helper method to find all Buildings taller than minimum and larger than minimum area given
*/
func (building *Building) FindAllBuildingsTallerAndWider(session *mgo.Session, h float64, a float64) ([]Building, error) {
	var buildings []Building
	err := getDBCollection(session).Find(bson.M{"$and": []bson.M{bson.M{"heightRoof": bson.M{"$gt": h}},
		bson.M{"shapeArea": bson.M{"$gt": a}}}}).All(&buildings)
	return buildings, err
}

/**
Helper method to find all demolished structures constructed in a given year
*/
func (building *Building) FindAllDemolishedStructruesByYear(session *mgo.Session, y int) ([]Building, error) {
	var buildings []Building
	err := getDBCollection(session).Find(bson.M{"$and": []bson.M{bson.M{"year": y},
		bson.M{"lastStatus": "Demolition"}}}).All(&buildings)
	return buildings, err
}

/**
Helper method to check if the struct fields are empty
*/
func isEmpty(k reflect.Kind, e reflect.Value) bool {

	empty := false

	switch k {
	case reflect.String:
		if e.String() == "" {
			empty = true
		}

	case reflect.Float32, reflect.Float64:
		if e.Float() == 0 {
			empty = true
		}

	case reflect.Int16, reflect.Int32, reflect.Int:
		if e.Int() == 0 {
			empty = true
		}
	}
	return empty
}
