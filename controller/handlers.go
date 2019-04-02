package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strconv"
	"time"
	"topos-backend-assignment/db"
	"topos-backend-assignment/models"
)

var Logger *log.Logger
var js []byte
var valid bool

type ResponseMessage struct {
	Status   int    `json:"status"`
	ErrorMsg string `json:"error"`
	Message  string `json:"message"`
}

func UseLogger(l *log.Logger) {
	Logger = l
}

/**
API URL - http://<host>:<port>/
Method	- GET
Params	- None
This API returns the count of all the BuildingFootPrint Data in the Database
*/
func AllBuildingsCount(writer http.ResponseWriter, request *http.Request) {
	var err error
	var bld models.Building
	var msg ResponseMessage
	var cnt int

	valid, js = IsTokenValid(request)
	if valid {
		// token is valid and no other error, proceed
		err, cnt = bld.GetAllBuildingsCount(db.MgoSession)
		if err != nil {
			Logger.Println("Error getting building count ", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error getting Building count"}
			js, err = json.Marshal(msg)
		} else {
			msg = ResponseMessage{Status: http.StatusOK, ErrorMsg: "", Message: "Total Building Footprints " + strconv.Itoa(cnt)}
			js, err = json.Marshal(msg)
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://<host>:<port>/buildingFootprints
Method	- GET
Params	- None
This API returns all the BuildingFootPrint Data from the Database
*/
func AllBuildingFootPrints(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	var err error
	var buildings []models.Building
	var msg ResponseMessage

	valid, js = IsTokenValid(request)
	if valid {
		// token is valid and no other error, proceed
		buildings, err = bld.GetAllBuildingFootPrints(db.MgoSession)
		if err != nil {
			Logger.Println("Error getting all building footprints...", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error getting FootPrint Data"}
			js, err = json.Marshal(msg)
		} else {
			if len(buildings) != 0 {
				js, err = json.Marshal(buildings)
			} else {
				Logger.Println("No data found")
				msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "Empty Dataset"}
				js, err = json.Marshal(msg)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://<host>:<port>/buildingFootprints/{id}
Method	- GET
Params	- id : Building footprint ID
This API returns the BuildingFootPrint Data for a specific ID
*/
func GetBuildingFootPrintsById(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	params := mux.Vars(request)
	var js []byte
	var err error
	var msg ResponseMessage

	valid, js = IsTokenValid(request)
	if valid {
		// token is valid and no other error, proceed
		bld, err = bld.FindBuildingFootPrintById(db.MgoSession, params["id"])
		if err != nil {
			Logger.Println("Error getting building footprint by ID...", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error getting Building FootPrint Data"}
			js, err = json.Marshal(msg)
		} else {
			if (models.Building{}) != bld {
				js, err = json.Marshal(bld)
			} else {
				Logger.Println("No data found")
				msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "Empty Dataset"}
				js, err = json.Marshal(msg)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://<host>:<port>/buildingFootprints/{id}
Method	- DELETE
Params	- id : Building footprint ID
This API deletes the BuildingFootPrint Data for a specific ID
*/
func DeleteBuildingFootPrints(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	params := mux.Vars(request)
	var js []byte
	var err error
	var msg ResponseMessage

	valid, js = IsTokenValid(request)
	if valid {
		// token is valid and no other error, proceed
		bld, err = bld.FindBuildingFootPrintById(db.MgoSession, params["id"])
		if err != nil {
			Logger.Println("Error deleting building footprint data...", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Deleting FootPrint Data"}
			js, err = json.Marshal(msg)
		} else {
			if (models.Building{}) != bld {
				if err = bld.DeleteBuildingFootPrint(db.MgoSession, bld); err != nil {
					Logger.Println("Error deleting building footprint data...", err)
					msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Deleting FootPrint Data"}
					js, err = json.Marshal(msg)
				} else {
					msg = ResponseMessage{Status: http.StatusOK, ErrorMsg: "None", Message: "BuildingFootPrint successfully deleted!"}
					js, err = json.Marshal(msg)
				}
			} else {
				msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "No Such BuildingFootPrint"}
				js, err = json.Marshal(msg)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://<host>:<port>/buildingFootprints/{id}
Method	- PUT
Params	- id : Building footprint ID
Request Body - JSON Object containing all the fields to be updated
This API updates the BuildingFootPrint Data for a specific ID
*/
func UpdateBuildingFootPrints(writer http.ResponseWriter, request *http.Request) {

	var bld models.Building
	var bldUpdate models.Building
	params := mux.Vars(request)
	var js []byte
	var err error
	var msg ResponseMessage

	valid, js = IsTokenValid(request)
	if valid {
		// token is valid and no other error, proceed
		bld, err = bld.FindBuildingFootPrintById(db.MgoSession, params["id"])
		if err != nil {
			Logger.Println("Error updating building footprint data...", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Updating FootPrint Data"}
			js, err = json.Marshal(msg)
		} else {
			if (models.Building{}) != bld {

				if err = json.NewDecoder(request.Body).Decode(&bldUpdate); err != nil {
					Logger.Println("Error parsing request...", err)
					msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error parsing request"}
					js, err = json.Marshal(msg)
				} else {
					if err = bld.UpdateBuildingFootPrint(db.MgoSession, bld, bldUpdate); err != nil {
						Logger.Println("Error updating building footprint data...", err)
						msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Updating FootPrint Data"}
						js, err = json.Marshal(msg)
					} else {
						msg = ResponseMessage{Status: http.StatusOK, ErrorMsg: "None", Message: "BuildingFootPrint successfully updated!"}
						js, err = json.Marshal(msg)
					}
				}
			} else {
				msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "No Such BuildingFootPrint"}
				js, err = json.Marshal(msg)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://<host>:<port>/buildingFootprints
Method	- POST
Params	- None
Request Body - JSON Object containing all the fields
This API deletes the BuildingFootPrint Data for a specific ID
*/
func AddBuildingFootPrints(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	var js []byte
	var err error
	var msg ResponseMessage
	defer request.Body.Close()

	valid, js = IsTokenValid(request)
	if valid {
		// token is valid and no other error, proceed
		if err = json.NewDecoder(request.Body).Decode(&bld); err != nil {
			Logger.Println("Error adding building footprint data...", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Creating FootPrint Data"}
			js, err = json.Marshal(msg)
		} else {
			bld.ID = bson.NewObjectId()
			bld.LastModified = time.Now()
			if err = bld.CreateBuildingFootPrint(db.MgoSession, bld); err != nil {
				Logger.Println("Error adding building footprint data...", err)
				msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Creating FootPrint Data"}
				js, err = json.Marshal(msg)
			} else {
				msg = ResponseMessage{Status: http.StatusCreated, ErrorMsg: "None", Message: "BuildingFootPrint successfully created!"}
				js, err = json.Marshal(msg)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://<host>:<port>/buildingFootprints/type/{bldType}
Method	- GET
Params	- bldType : feat code of the building type
This API returns all the Buildings by Type
*/
func GetBuildingsByType(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	params := mux.Vars(request)
	var js []byte
	var err error
	var buildings []models.Building
	var msg ResponseMessage
	var bldngType int

	valid, js = IsTokenValid(request)
	if valid {
		// token is valid and no other error, proceed
		bldngType, err = strconv.Atoi(params["bldType"])
		buildings, err = bld.FindAllBuildingsByType(db.MgoSession, bldngType)
		if err != nil {
			Logger.Println("Error getting buildings by type...", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error getting Buildings by type"}
			js, err = json.Marshal(msg)
		} else {
			if len(buildings) != 0 {
				js, err = json.Marshal(buildings)
			} else {
				msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "Empty Dataset"}
				js, err = json.Marshal(msg)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://<host>:<port>/buildingFootprints/buildingHeightAndArea/{minHeight}/{minArea}
Method	- GET
Params	- minHeight : the height with reference to which taller buildings are to be returned
Params	- minArea : the area with reference to which larger area buildings are to be returned
This API returns all the Buildings taller and larger than the given values
*/
func GetTallAndWideBuildings(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	params := mux.Vars(request)
	var js []byte
	var err error
	var buildings []models.Building
	var msg ResponseMessage
	var height float64
	var area float64

	valid, js = IsTokenValid(request)
	if valid {
		// token is valid and no other error, proceed
		height, err = strconv.ParseFloat(params["minHeight"], 32)
		area, err = strconv.ParseFloat(params["minArea"], 32)
		buildings, err = bld.FindAllBuildingsTallerAndWider(db.MgoSession, height, area)
		if err != nil {
			Logger.Println("Error getting taller and wider buildings...", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error getting taller and wider buildings"}
			js, err = json.Marshal(msg)
		} else {
			if len(buildings) != 0 {
				js, err = json.Marshal(buildings)
			} else {
				msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "No such data present"}
				js, err = json.Marshal(msg)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
API URL - http://<host>:<port>/buildingFootprints/demolishedStructuresByConstructedYear/{year}
Method	- GET
Params	- year : Year for the demolished structures to be returned
This API returns all the demolished structures in a given year
*/
func GetAllDemolishedStructuresByYear(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	params := mux.Vars(request)
	var js []byte
	var err error
	var buildings []models.Building
	var msg ResponseMessage
	var y int

	valid, js = IsTokenValid(request)
	if valid {
		// token is valid and no other error, proceed
		y, err = strconv.Atoi(params["year"])
		buildings, err = bld.FindAllDemolishedStructruesByYear(db.MgoSession, y)
		if err != nil {
			Logger.Println("Error getting demolished structures...", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error getting demolished structures"}
			js, err = json.Marshal(msg)
		} else {
			if len(buildings) != 0 {
				js, err = json.Marshal(buildings)
			} else {
				msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "No such data present"}
				js, err = json.Marshal(msg)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}
