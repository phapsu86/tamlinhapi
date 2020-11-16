package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
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

type ItemSearch struct {
	ReligionID int    `json:"religion_id"`
	Keyword    string `json:"keyword"`
	LocID      string `json:"loc_id"`
	Page       int    `json:"page"`
}

type ParamItemFollow struct {
	ID    uint64 `json:"id"`
	Value int    `json:"value"`
}

func (server *Server) CreateReligionItem(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	item := models.ReligionItem{}
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
	itemCreated, err := item.SaveReligionItem(server.DB)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, itemCreated.ID))
	responses.JSON(w, http.StatusCreated, itemCreated)
}

func (server *Server) GetReligionItems(w http.ResponseWriter, r *http.Request) {

	model := models.ReligionItem{}

	items, err := model.FindAllReligionItems(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
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
	//	responses.JSON(w, http.StatusOK, items)
}

func (server *Server) GetReligionItemByReligionID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	page, err := strconv.ParseUint(vars["page"], 10, 64)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	model := models.ReligionItem{}

	items, err := model.FindReligionItemtByReligionID(server.DB, pid, page)
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
	//result := formatresult.ReturnReligionItem(items)
	//responses.JSON(w, http.StatusOK, result)
}

func (server *Server) GetReligionItemDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	model := models.ReligionItem{}

	item, err := model.FindReligionItemDetails(server.DB, pid)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	linkImg := server.getLink(item.Image, "tamlinh")

	item.Image = linkImg

	var rs interface{}
	rs = item

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)
	//	result := formatresult.ReturnReligionItemDetail(itemReceived)
	//	responses.JSON(w, http.StatusOK, result)
}

//Search item by keyword

func (server *Server) SearchItemByNames(w http.ResponseWriter, r *http.Request) {

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	modelSearch := ItemSearch{}
	err = json.Unmarshal(body, &modelSearch)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	locID := modelSearch.LocID
	Keyword := modelSearch.Keyword
	ReligionID := modelSearch.ReligionID

	page := modelSearch.Page

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	model := models.ReligionItem{}

	items, err := model.SearchItemByName(server.DB, ReligionID, Keyword, locID, page)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.Image, "tamlinh")

			items[i].Image = linkImg

		}
	}

	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

	//result := formatresult.ReturnReligionItem(itemReceived)
	//responses.JSON(w, http.StatusOK, result)
}

//Search item by keyword

func (server *Server) CreatItemFollow(w http.ResponseWriter, r *http.Request) {

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data

	modelFL := models.ReligionItemFollow{}
	pr := ParamItemFollow{}
	err = json.Unmarshal(body, &pr)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	mdReligionItem := models.ReligionItem{}

	_, err = mdReligionItem.FindReligionItemDetails(server.DB, pr.ID)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("RELIGION_ITEM_DOES_NOT_EXIST"))
		responses.JSON(w, http.StatusOK, err)
		return
	}

	if pr.Value == 1 {
		_, err = modelFL.FindUserReligionItemFollowID(server.DB, uid, pr.ID)
		if err == nil {
			err := formaterror.ReturnErr(errors.New("ITEM_EXIST"))
			responses.JSON(w, http.StatusOK, err)
			return
		}
	}

	modelFL.Prepare()
	modelFL.ItemID = pr.ID
	modelFL.UserID = uid
	err = modelFL.SaveUserItemFollow(server.DB, uid, pr.ID, pr.Value)

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	result := formaterror.ReturnSuccess()
	responses.JSON(w, http.StatusOK, result)
}

func (server *Server) GetReligionItemFollow(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}

	// Start processing the request data

	modelFL := models.ReligionItemFollow{}
	items, err := modelFL.FindAllReligionItemFollow(server.DB, uid)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)
	//	responses.JSON(w, http.StatusOK, rs)
}

func (server *Server) GetInfoItemFollowByID(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}

	vars := mux.Vars(r)
	itemID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	// Start processing the request data
	modelFL := models.ReligionItemFollow{}
	_, err = modelFL.GetReligionItemFollowByID(server.DB, uid, itemID)
	var rs interface{}
	if err != nil {
		rs = 0
		result := formatresult.ReturnGlobalArray(rs)
		responses.JSON(w, http.StatusOK, result)
		return
	}

	rs = 1
	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)
	//	responses.JSON(w, http.StatusOK, rs)
}

func (server *Server) CheckUserOfItem(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}
	vars := mux.Vars(r)
	itemID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	// Start processing the request data
	modelUOI := models.UserOfItem{}
	item, err := modelUOI.FindUserOfItemByID(server.DB, uid, itemID)
	var rs interface{}
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return

	}

	rs = item
	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

}
