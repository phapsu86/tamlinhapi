package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/victorsteven/fullstack/api/auth"
	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"
	"github.com/victorsteven/fullstack/api/utils/formaterror"
	"github.com/victorsteven/fullstack/api/utils/formatresult"
)

type ResponseResult struct {
	Status int        `json:"status"`
	Data   ResultLSFC `json:"data"`
}

type ResultLSFC struct {
	Like    int64 `json:"like"`
	Share   int64 `json:"share"`
	Join    int64 `json:"join"`
	Follow  int64 `json:"follow"`
	Comment int64 `json:"comment"`
}


type ParamLSFJ struct {
	Value int        `json:"value"`
	ActionType  int `json:"action_type"`
	ID  uint64 `json:"id"`
}


func (server *Server) CreateReligionEventLsfj(w http.ResponseWriter, r *http.Request) {

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	item := ParamLSFJ{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	modelLSFJ := models.ReligionEventLsfj{}
	modelLSFJ.Prepare()
	modelLSFJ.UserID = uid

	 err = modelLSFJ.SaveReligionEventLsfj(server.DB,uid,item.ID,item.Value,item.ActionType)
	if err != nil {
		formattedError := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusInternalServerError, formattedError)
		return
	}
	rs := formaterror.ReturnSuccess()
	responses.JSON(w, http.StatusCreated, rs)
}

// func (server *Server) GetReligionEventLsfjs(w http.ResponseWriter, r *http.Request) {

// 	model := models.ReligionEventLsfj{}

// 	items, err := model.FindAllReligionEventLsfjs(server.DB)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	responses.JSON(w, http.StatusOK, items)
// }

func (server *Server) GetTotalEventLsfjByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	model := models.ReligionEventLsfj{}

	j, s, f, err := model.GetTotalEventLsfjs(server.DB, pid)

	//var item = [3]int64{j, s, f}
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	rs := ResultLSFC{}

	rs.Join = j
	rs.Share = s
	rs.Follow = f
	resp := ResponseResult{}
	resp.Status = 200
	resp.Data = rs

	responses.JSON(w, http.StatusOK, resp)
}

func (server *Server) UpdateReligionEventLsfj(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	model := models.ReligionEventLsfj{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&model).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	// If a user attempt to update a post not belonging to him
	if uid == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	modeltUpdate := models.ReligionEventLsfj{}
	err = json.Unmarshal(body, &modeltUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	modeltUpdate.Prepare()
	err = modeltUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	modeltUpdate.ID = model.ID //this is important to tell the model the post id to update, the other update field are set above

	modelUpdated, err := modeltUpdate.UpdateAReligionEventLsfj(server.DB)

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	responses.JSON(w, http.StatusOK, modelUpdated)
}

func (server *Server) DeleteReligionEventLsfj(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	model := models.ReligionEventLsfj{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&model).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = model.DeleteAReligionEventLsfj(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}

func (server *Server) GetInfoEventLSFJByID(w http.ResponseWriter, r *http.Request) {
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

	lsfjModel := models.ReligionEventLsfj{}

	item, err := lsfjModel.FindReligionEventLsfjByID(server.DB, uid, pid)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	var rs interface{}
	rs = item

	result := formatresult.ReturnGlobalArray(rs)
	
	responses.JSON(w, http.StatusOK, result)
}




func (server *Server) GetEventLSFJFollow(w http.ResponseWriter, r *http.Request) {
	//resultFail := ResultFail{}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}

	vars := mux.Vars(r)

	// Is a valid post id given to us?
	page, err := strconv.ParseUint(vars["page"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	lsfjModel := models.ReligionEventLsfj{}
	items, err := lsfjModel.FindEventLsfjFollow(server.DB, uid, page)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}


	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.EventDetails.Image, "tamlinh")
			//fmt.Printf("dfdfdfd %v", linkImg)
			items[i].EventDetails.Image = linkImg
			fmt.Printf("dfdfdfd %v", v)

		}
	}


	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

	//result := formatresult.ReturnEventLSFJ(lsfjData)
	
	//responses.JSON(w, http.StatusOK, result)
}