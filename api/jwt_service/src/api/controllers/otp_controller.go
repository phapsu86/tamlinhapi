package controllers

import (
	"crypto/rand"
	"io"
	"net/http"

	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/mux"

	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"

	"github.com/victorsteven/fullstack/api/utils/formaterror"
)

//Chưc năng get test
var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func (server *Server) GetOTP(w http.ResponseWriter, r *http.Request) {

	post := models.Otp{}
	vars := mux.Vars(r)
	phone := vars["phone"]

	posts, err := post.FindOtpPhone(server.DB, phone)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

func (server *Server) SendOTP(w http.ResponseWriter, r *http.Request) {

	otpMD := models.Otp{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &otpMD)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	otpStr := EncodeToString(6)
	otpMD.Code = otpStr
	otpMD.Prepare()
	err = otpMD.CreateOtp(server.DB)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	result := formaterror.ReturnSuccess()
	responses.JSON(w, http.StatusOK, result)
}

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
