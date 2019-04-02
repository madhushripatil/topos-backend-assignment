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

type Borough struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	BoroughCode int64         `bson:"boroughCode" json:"boroughCode"`
	BoroughName string        `bson:"boroughName" json:"boroughName"`
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
		BoroughCollectionName := os.Getenv("borough_collection_name")
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

		c := mongoSession.DB(DbName).C(BoroughCollectionName)
		file, err := os.Open("borough.csv")
		if err != nil {
			panic(err)
		}

		defer file.Close()
		reader := csv.NewReader(file)
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

				boroughCode, err := strconv.ParseInt(record[0], 10, 64)
				if err != nil {
					fmt.Println("Error in parsing value", err)
				}

				document := &Borough{
					ID:          bson.NewObjectId(),
					BoroughCode: boroughCode,
					BoroughName: record[1],
				}

				err = c.Insert(document)

				if err != nil {
					panic(err)
				}
			}
			fmt.Println("Building FeatType data imported successfully!")
		}
	}
}
