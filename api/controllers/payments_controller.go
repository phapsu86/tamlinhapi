package controllers

import (
	"encoding/json"
	"errors"

	//	"fmt"
	"io/ioutil"
	"net/http"

	//"strconv"

	//"github.com/gorilla/mux"
	"github.com/phapsu86/tamlinhapi/api/auth"
	"github.com/phapsu86/tamlinhapi/api/models"

	"github.com/phapsu86/tamlinhapi/api/responses"
)

type ParamTopupPoint struct {
	PaymentType string `json:"payment_type"`
	Amount      int    `json:"status"`
	Page        int    `json:"page"`
}

func (server *Server) CreatePaymentTopUpPoint(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	item := ParamTopupPoint{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil || uid == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//Check mo
	//playLoad := payments.Payload{}
	//stt, err := payments.paymentProcess()
	//	mdOrder := models.Order{}
	//	mdTransHis := models.TransactionHistory{}
	//	mdOrderDetail := models.OrderDetail{}

}

func (server *Server) CreateTransactions(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	item := models.ReligionList{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil || uid == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//Tạo oreder

	//Tạo transaction

	//Tra ve mã hóa đơn và refid

}
