package controller

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"topos-backend-assignment/db"
	"topos-backend-assignment/models"
)

/**
API URL - http://<host>:<port>/transform/archiveDemolishedStructures
Method	- POST
Params	- None
Header	- Authorization : <JWT>
This API archives all the demolished structures
*/
func ArchiveAllDemolishedStructures(writer http.ResponseWriter, request *http.Request) {
	var bld models.Building
	var js []byte
	var err error
	var buildings []models.Building
	var demolishedStruct models.Demolished
	var msg ResponseMessage
	flag := true

	valid, js = IsTokenValid(request)
	if valid {
		// token is valid and no other error, proceed
		buildings, err = bld.FindAllDemolishedStructures(db.MgoSession)
		if err != nil {
			Logger.Println("Error getting demolished structures...", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error getting demolished structures"}
			js, err = json.Marshal(msg)
		} else {
			length := len(buildings)
			if length != 0 {

				for i := 0; i < length; i++ {
					//archive demolished structure
					demolishedStruct = models.Demolished(buildings[i])
					demolishedStruct.ID = bson.NewObjectId()
					if err = demolishedStruct.CreateDemolishedStructure(db.MgoSession, demolishedStruct); err != nil {
						Logger.Println("Error archiving demolished structures", err)
						msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error archiving demolished structures"}
						js, err = json.Marshal(msg)
						flag = false
						break
					} else {
						//delete archived structure from source
						if err = bld.DeleteBuildingFootPrint(db.MgoSession, buildings[i]); err != nil {
							Logger.Println("Error archiving demolished structures", err)
							msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error archiving demolished structures"}
							js, err = json.Marshal(msg)
							flag = false
							break
						}
					}
				}
				if flag {
					msg = ResponseMessage{Status: http.StatusOK, ErrorMsg: "", Message: "Demolished structures archived successfully!"}
					js, err = json.Marshal(msg)
				}
			} else {
				msg = ResponseMessage{Status: http.StatusNoContent, ErrorMsg: "", Message: "No data to be archived!"}
				js, err = json.Marshal(msg)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}
