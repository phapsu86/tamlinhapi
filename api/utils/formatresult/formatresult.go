package formatresult

import (
	//"errors"
	//"strings"
	"github.com/phapsu86/tamlinhapi/api/models"
)

type ResultGlobalArray struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func ReturnGlobalArray(data interface{}) ResultGlobalArray {
	model := ResultGlobalArray{}
	model.Status = 200
	model.Data = data
	return model
}



type ResultUser struct {
	Status int         `json:"status"`
	Data   models.User `json:"data"`
}

func ReturnUser(data *models.User) ResultUser {
	model := ResultUser{}
	model.Status = 200
	model.Data = *data
	return model
}

type ResultUserRemind struct {
	Status int                 `json:"status"`
	Data   []models.UserRemind `json:"data"`
}

func ReturnUserRemind(data *[]models.UserRemind) ResultUserRemind {
	model := ResultUserRemind{}
	model.Status = 200
	model.Data = *data
	return model
}

//Format Post

type ResultPost struct {
	Status int           `json:"status"`
	Data   []models.Post `json:"data"`
}

func ReturnPost(data *[]models.Post) ResultPost {
	model := ResultPost{}
	model.Status = 200
	model.Data = *data
	return model
}

type ResultPostDetail struct {
	Status int         `json:"status"`
	Data   models.Post `json:"data"`
}

func ReturnPostDetail(data *models.Post) ResultPostDetail {
	model := ResultPostDetail{}
	model.Status = 200
	model.Data = *data
	return model
}




type ResultPostComment struct {
	Status int           `json:"status"`
	Data   []models.PostComment `json:"data"`
}

func ReturnPostComment(data *[]models.PostComment) ResultPostComment {
	model := ResultPostComment{}
	model.Status = 200
	model.Data = *data
	return model
}




//==================Event=========================

type ResultEvents struct {
	Status int                    `json:"status"`
	Data   []models.ReligionEvent `json:"data"`
}

func ReturnEvents(data *[]models.ReligionEvent) ResultEvents {
	model := ResultEvents{}
	model.Status = 200
	model.Data = *data
	return model
}

type ResultEventDetail struct {
	Status int                  `json:"status"`
	Data   models.ReligionEvent `json:"data"`
}

func ReturnEventDetail(data *models.ReligionEvent) ResultEventDetail {
	model := ResultEventDetail{}
	model.Status = 200
	model.Data = *data
	return model
}

//Chùa
type ResultReligionItem struct {
	Status int                   `json:"status"`
	Data   []models.ReligionItem `json:"data"`
}

func ReturnReligionItem(data *[]models.ReligionItem) ResultReligionItem {
	model := ResultReligionItem{}
	model.Status = 200
	model.Data = *data
	return model
}

type ResultReligionItemDetail struct {
	Status int                 `json:"status"`
	Data   models.ReligionItem `json:"data"`
}

func ReturnReligionItemDetail(data *models.ReligionItem) ResultReligionItemDetail {
	model := ResultReligionItemDetail{}
	model.Status = 200
	model.Data = *data
	return model
}

//Công Duc
type ResultMerit struct {
	Status int                `json:"status"`
	Data   []models.MeritList `json:"data"`
}

func ReturnResultMerit(data *[]models.MeritList) ResultMerit {
	model := ResultMerit{}
	model.Status = 200
	model.Data = *data
	return model
}

type ResultResultMeritDetail struct {
	Status int              `json:"status"`
	Data   models.MeritList `json:"data"`
}

func ReturnResultMeritDetail(data *models.MeritList) ResultResultMeritDetail {
	model := ResultResultMeritDetail{}
	model.Status = 200
	model.Data = *data
	return model
}

//Vat pham  dc ban ch
type ResultOfferingItemSell struct {
	Status int                       `json:"status"`
	Data   []models.OfferingItemSell `json:"data"`
}

func ReturnOfferingItemSell(data *[]models.OfferingItemSell) ResultOfferingItemSell {
	model := ResultOfferingItemSell{}
	model.Status = 200
	model.Data = *data
	return model
}


// LSFS cho event 
type ResultEventLSFJ struct {
	Status int                    `json:"status"`
	Data   models.ReligionEventLsfj `json:"data"`
}

func ReturnEventLSFJ (data *models.ReligionEventLsfj) ResultEventLSFJ {
	model := ResultEventLSFJ{}
	model.Status = 200
	model.Data = *data
	return model
}


// LSFS cho POST
type ResultPostLSFJ struct {
	Status int                    `json:"status"`
	Data   models.PostLsfc `json:"data"`
}

func ReturnPostLSFJ (data *models.PostLsfc) ResultPostLSFJ {
	model := ResultPostLSFJ{}
	model.Status = 200
	model.Data = *data
	return model
}



