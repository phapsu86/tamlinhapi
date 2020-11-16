package controllers

import (
	"encoding/json"
	"errors"

	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/phapsu86/tamlinhapi/api/auth"
	"github.com/phapsu86/tamlinhapi/api/models"
	"github.com/phapsu86/tamlinhapi/api/responses"
	"github.com/phapsu86/tamlinhapi/api/utils/formaterror"
	"github.com/phapsu86/tamlinhapi/api/utils/formatresult"
)

func (server *Server) CreatePostLsfc(w http.ResponseWriter, r *http.Request) {

	uid, err := auth.ExtractTokenID(r)
	if err != nil || uid == 0 {
		//responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	item := ParamLSFJ{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	modelLSFJ := models.PostLsfc{}
	modelLSFJ.Prepare()
	modelLSFJ.UserID = uid

	err = modelLSFJ.SavePostLsfc(server.DB, uid, item.ID, item.Value, item.ActionType)
	if err != nil {
		formattedError := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusInternalServerError, formattedError)
		return
	}
	rs := formaterror.ReturnSuccess()
	responses.JSON(w, http.StatusCreated, rs)
}

// func (server *Server) GetPostLsfcs(w http.ResponseWriter, r *http.Request) {

// 	model := models.PostLsfc{}

// 	items, err := model.FindAllPostLsfcs(server.DB)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	responses.JSON(w, http.StatusOK, items)
// }

func (server *Server) GetTotalPostLsfjByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	model := models.PostLsfc{}

	c, s, f, l, err := model.GetTotalPostLsfjc(server.DB, pid)

	//var item = [4]int64{c, s, f, l}
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	rs := ResultLSFC{}

	rs.Comment = c
	rs.Share = s
	rs.Follow = f
	rs.Like = l
	resp := ResponseResult{}
	resp.Status = 200
	resp.Data = rs

	responses.JSON(w, http.StatusOK, resp)
	//responses.JSON(w, http.StatusOK, item)
}

func (server *Server) GetPostLSFJByUser(w http.ResponseWriter, r *http.Request) {
	//resultFail := ResultFail{}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}

	vars := mux.Vars(r)

	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	lsfjModel := models.PostLsfc{}

	lsfjData, err := lsfjModel.FindReligionPostLsfjByID(server.DB, uid, pid)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	result := formatresult.ReturnPostLSFJ(lsfjData)

	responses.JSON(w, http.StatusOK, result)
}

func (server *Server) GetAllPostLSFJByUser(w http.ResponseWriter, r *http.Request) {
	//resultFail := ResultFail{}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}

//	vars := mux.Vars(r)

	// Is a valid post id given to us?
//	page, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	lsfjModel := models.PostLsfc{}

	items, err := lsfjModel.FindAllPostLsfjByUser(server.DB, uid, 0)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	var rs interface{}
	rs = items
	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

}
