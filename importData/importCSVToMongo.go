package main

import (
	"encoding/csv"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Building struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	Bin          int64         `bson:"bin" json:"bin"`
	BoroughCode  int           `bson:"boroughCode" json:"boroughCode"`
	BuildingCode int           `bson:"buildingCode" json:"buildingCode"`
	DoittID      int64         `bson:"doittID" json:"doittID"`
	Name         string        `bson:"name" json:"name"`
	Year         int64         `bson:"year" json:"year"`
	LastStatus   string        `bson:"lastStatus" json:"lastStatus"`
	FeatCode     int64         `bson:"featCode" json:"featCode"`
	HeightRoof   float64       `bson:"heightRoof" json:"heightRoof"`
	GroundLevel  int64         `bson:"groundLevel" json:"groundLevel"`
	ShapeArea    float64       `bson:"shapeArea" json:"shapeArea"`
	ShapeLength  float64       `bson:"shapeLength" json:"shapeLength"`
	LastModified time.Time     `bson:"lastModified" json:"lastModified"`
	GeomSource   string        `bson:"geomSource" json:"geomSource"`
	BaseBBL      int           `bson:"baseBBL" json:"baseBBL"`
	MPlutoBBL    int           `bson:"mplutoBBL" json:"mplutoBBL"`
}

func main() {

	log.Println("Loading the configuration file...")
	envFile := fmt.Sprintf("%s.env", os.Getenv("BUILDING_ENV"))
	e := godotenv.Load(envFile)

	if e != nil {
		// close the application if the configuration file open fails
		log.Println("Error loading the config file, Closing the application.")
		panic(e)
	} else {

		DbName := os.Getenv("db_name")
		BuildingCollectionName := os.Getenv("bldng_collection_name")
		DBHost := os.Getenv("db_host")
		DBTimeout := os.Getenv("db_timeout")
		DBUsername := os.Getenv("db_username")
		DBPassword := os.Getenv("db_pass")

		timeout, err := strconv.Atoi(DBTimeout)

		mongoDBDialInfo := &mgo.DialInfo{
			Addrs:    []string{DBHost},
			Timeout:  time.Duration(timeout) * time.Second,
			Database: DbName,
			Username: DBUsername,
			Password: DBPassword,
		}

		mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
		if err != nil {
			log.Println("CreateSession: %s\n", err)
			panic(err)
		} else {
			log.Println("Created DB Session. Connected to:", DbName)
		}

		defer mongoSession.Close()
		mongoSession.SetMode(mgo.Monotonic, true)

		c := mongoSession.DB(DbName).C(BuildingCollectionName)
		file, err := os.Open("data.csv")
		if err != nil {
			panic(err)
		}

		defer file.Close()
		reader := csv.NewReader(file)
		dateFormat := "01/02/2006 15:04:05 PM +0000"
		var str []string
		str, err = reader.Read()
		if err != nil {
			// close program if error occurs reading file
			panic(err)
		} else {
			fmt.Println(str)
			// load the csv file into database
			for {
				record, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					panic(err)
				}

				parsedLastModified, err := time.Parse(dateFormat, record[4])
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				bin, err := strconv.ParseInt(record[1], 10, 64)
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				boroughCode, err := strconv.Atoi(record[1][0:1])
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				buildingCode, err := strconv.Atoi(record[1][1:])
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				year, err := strconv.ParseInt(record[2], 10, 64)
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				doittId, err := strconv.ParseInt(record[6], 10, 64)
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				heightRoof, err := strconv.ParseFloat(record[7], 64)
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				featCode, err := strconv.ParseInt(record[8], 10, 64)
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				groundLevel, err := strconv.ParseInt(record[9], 10, 64)
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				shapeArea, err := strconv.ParseFloat(record[10], 64)
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				shapeLength, err := strconv.ParseFloat(record[11], 64)
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				baseBBL, err := strconv.Atoi(record[12])
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				mplutoBBL, err := strconv.Atoi(record[13])
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				document := &Building{
					ID:           bson.NewObjectId(),
					Bin:          bin,
					BoroughCode:  boroughCode,
					BuildingCode: buildingCode,
					Year:         year,
					Name:         record[3],
					LastModified: parsedLastModified,
					DoittID:      doittId,
					HeightRoof:   heightRoof,
					FeatCode:     featCode,
					GroundLevel:  groundLevel,
					ShapeArea:    shapeArea,
					ShapeLength:  shapeLength,
					BaseBBL:      baseBBL,

					MPlutoBBL:  mplutoBBL,
					GeomSource: record[14],
				}

				err = c.Insert(document)

				if err != nil {
					panic(err)
				}
			}
			fmt.Println("Data imported successfully!")
		}
	}
}
