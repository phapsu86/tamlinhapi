package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	//"strconv"

	//"github.com/gorilla/mux"
	"github.com/phapsu86/tamlinhapi/api/auth"
	"github.com/phapsu86/tamlinhapi/api/models"
	"github.com/phapsu86/tamlinhapi/api/responses"
	"github.com/phapsu86/tamlinhapi/api/utils/formaterror"
)


type ResultPoint struct { 

	Status   int  `json:"status"` 
	Data int64  `json:"data"` 
	
}

func (server *Server) CreatePoint(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	Point := models.Point{}
	err = json.Unmarshal(body, &Point)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	Point.Prepare()
	err = Point.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }

	PointCreated, err := Point.SavePoint(server.DB)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
		
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, PointCreated.ID))
	responses.JSON(w, http.StatusCreated, PointCreated)
}



func (server *Server) GetPoint(w http.ResponseWriter, r *http.Request) {
	//resultFail := ResultFail{}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}
	point := models.Point{}

	PointReceived, err := point.FindPointByID(server.DB, uid)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	result:= ResultPoint{}
	result.Status = http.StatusOK
	result.Data = PointReceived
	responses.JSON(w, http.StatusOK, result)
}

