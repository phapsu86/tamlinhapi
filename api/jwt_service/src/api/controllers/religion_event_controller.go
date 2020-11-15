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

type EventSearch struct {
	ReligionID int    `json:"religion_id"`
	Keyword    string `json:"keyword"`
	LocID      string    `json:"loc_id"`
	EventType  int    `json:"event_type"`
	Page       int    `json:"page"`
}

type EventSearchByDate struct {
	ReligionID int `json:"religion_id"`
	//Keyword   string  `json:"keyword"`
	LocID     string    `json:"loc_id"`
	EventType int    `json:"event_type"`
	FromDate  string `json:"from_date"`
	ToDate    string `json:"to_date"`
	Page      int    `json:"page"`
}

func (server *Server) CreateReligionEvent(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	item := models.ReligionEvent{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	item.Prepare()
	err = item.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil || uid == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// if uid != item.AuthorID {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }
	itemCreated, err := item.SaveReligionEvent(server.DB)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, itemCreated.ID))
	responses.JSON(w, http.StatusCreated, itemCreated)
}

func (server *Server) GetReligionEvents(w http.ResponseWriter, r *http.Request) {

	model := models.ReligionEvent{}

	items, err := model.FindAllReligionEvents(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, items)
}

// Danh sách sự kiên
func (server *Server) GetReligionEventByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	model := models.ReligionEvent{}

	page, err := strconv.ParseUint(vars["page"], 10, 64)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}

	items, err := model.FindReligionEventByID(server.DB, pid, page)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}


	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.Image, "tamlinh")
			//fmt.Printf("dfdfdfd %v", linkImg)
			items[i].Image = linkImg
			fmt.Printf("dfdfdfd %v", v)

		}
	}

	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)


	//result := formatresult.ReturnEvents(itemReceived)
	//responses.JSON(w, http.StatusOK, result)
}

func (server *Server) GetReligionEventDetails(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	model := models.ReligionEvent{}

	item, err := model.FindReligionEventDetail(server.DB, pid)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	linkImg := server.getLink(item.Image, "tamlinh")
	//fmt.Printf("dfdfdfd %v", linkImg)
	item.Image = linkImg

	var rs interface{}
	rs = item

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)
}

func (server *Server) GetReligionItemEventByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	t, err := strconv.ParseUint(vars["type"], 10, 64)
	page, err := strconv.ParseUint(vars["page"], 10, 64)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	model := models.ReligionEvent{}

	items, err := model.FindReligionItemEventByID(server.DB, pid, t, page)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.Image, "tamlinh")
			//fmt.Printf("dfdfdfd %v", linkImg)
			items[i].Image = linkImg
			fmt.Printf("dfdfdfd %v", v)

		}
	}

	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

	// result := formatresult.ReturnEvents(items)
	// responses.JSON(w, http.StatusOK, result)
}

func (server *Server) SearchEventByNames(w http.ResponseWriter, r *http.Request) {

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	modelSearch := EventSearch{}
	err = json.Unmarshal(body, &modelSearch)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	locID := modelSearch.LocID
	Keyword := modelSearch.Keyword
	EventType := modelSearch.EventType
	ReligionID := modelSearch.ReligionID

	page := modelSearch.Page

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	model := models.ReligionEvent{}

	items, err := model.SearchEventByName(server.DB, ReligionID, Keyword, EventType, locID, page)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}


	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.Image, "tamlinh")
			//fmt.Printf("dfdfdfd %v", linkImg)
			items[i].Image = linkImg
		//	fmt.Printf("dfdfdfd %v", v)

		}
	}



	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

}

//Where("created_at BETWEEN ? AND ?", lastWeek, today)

func (server *Server) SearchEventByDate(w http.ResponseWriter, r *http.Request) {

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	modelSearch := EventSearchByDate{}
	err = json.Unmarshal(body, &modelSearch)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	locID := modelSearch.LocID

	EventType := modelSearch.EventType
	ReligionID := modelSearch.ReligionID
	FromDate := modelSearch.FromDate
	ToDate := modelSearch.ToDate

	page := modelSearch.Page

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	model := models.ReligionEvent{}

	items, err := model.SearchEventByDate(server.DB, ReligionID, FromDate, ToDate, EventType, locID, page)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}


	
	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.Image, "tamlinh")
			//fmt.Printf("dfdfdfd %v", linkImg)
			items[i].Image = linkImg
		//	fmt.Printf("dfdfdfd %v", v)

		}
	}



	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

	//result := formatresult.ReturnEvents(itemReceived)
	//responses.JSON(w, http.StatusOK, result)
}

func (server *Server) UpdateReligionEvent(w http.ResponseWriter, r *http.Request) {

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
	model := models.ReligionEvent{}
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
	modeltUpdate := models.ReligionEvent{}
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

	modelUpdated, err := modeltUpdate.UpdateAReligionEvent(server.DB)

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	responses.JSON(w, http.StatusOK, modelUpdated)
}

func (server *Server) DeleteReligionEvent(w http.ResponseWriter, r *http.Request) {

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
	model := models.ReligionEvent{}
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
	_, err = model.DeleteAReligionEvent(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
