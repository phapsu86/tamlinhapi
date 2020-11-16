package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/phapsu86/tamlinhapi/api/auth"
	"github.com/phapsu86/tamlinhapi/api/models"
	"github.com/phapsu86/tamlinhapi/api/responses"
	"github.com/phapsu86/tamlinhapi/api/utils/formaterror"
	"github.com/phapsu86/tamlinhapi/api/utils/formatresult"
)

type RemindParam struct {
	ID         uint64 `json:"id"`
	RemindType int    `json:"remind_type"`
}

func (server *Server) CreateRemind(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Format errors"))
		responses.JSON(w, http.StatusOK, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}
	item := RemindParam{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("FORMAT_PARAM_ERROR"))
		responses.JSON(w, http.StatusOK, err)
		return
	}
	remindModel := models.UserRemind{}
	log.Printf("Successfully created %s\n", item.RemindType)
	if item.RemindType == 1 {
		// Insert UserRemind

		itemExist, _ := remindModel.FindUserRemindByEventID(server.DB, uid, item.ID)
		if itemExist.ID != 0 {
			err := formaterror.ReturnErr(errors.New("REMIND_EXIST"))
			responses.JSON(w, http.StatusOK, err)
			return

		}

		remindModel.Prepare()
		err = remindModel.Validate()
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		eventModel := models.ReligionEvent{}

		evData, err := eventModel.FindReligionEventDetail(server.DB, item.ID)
		if err != nil {
			err := formaterror.ReturnErr(errors.New("EVENT_NOT_FOUND"))
			responses.JSON(w, http.StatusOK, err)
			return
		}

		remindModel.ObjectID = item.ID
		remindModel.ObjectType = 1
		remindModel.BeginTime = evData.BeginDate
		remindModel.EndTime = evData.EndDate
		remindModel.UserID = uid
		remindModel.Status = 1

		itemCreated, err := remindModel.SaveUserRemind(server.DB)

		if err != nil {
			err := formaterror.ReturnErr(err)
			responses.JSON(w, http.StatusOK, err)
			return
		}

		if itemCreated.ID != 0 {
			sucess := formaterror.ReturnSuccess()
			responses.JSON(w, http.StatusOK, sucess)
		}

	} else {
		// xoa event

		_, err := remindModel.DeleteAUserRemind(server.DB, item.ID, uid)
		if err != nil {
			err := formaterror.ReturnErr(err)
			responses.JSON(w, http.StatusOK, err)
			return
		}
		sucess := formaterror.ReturnSuccess()
		responses.JSON(w, http.StatusOK, sucess)

	}

}

func (server *Server) GetListUserRemind(w http.ResponseWriter, r *http.Request) {

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		errResult := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, errResult)
		return
	}

	modelRemind := models.UserRemind{}
	items, err := modelRemind.FindAllUserReminds(server.DB, uint64(tokenID))
	//Get link for avartar
	if err != nil {
		responses.JSON(w, http.StatusOK, formaterror.ReturnErr(err))
		return
	}

	
	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.EventDetail.Image, "tamlinh")
			items[i].EventDetail.Image = linkImg
		}
	}



	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)


	// result := formatresult.ReturnUserRemind(dataRemind)

	// responses.JSON(w, http.StatusOK, result)
}
