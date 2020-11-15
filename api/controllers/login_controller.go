package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"

	"github.com/victorsteven/fullstack/api/auth"
	"github.com/victorsteven/fullstack/api/models"
	"github.com/victorsteven/fullstack/api/responses"
	"github.com/victorsteven/fullstack/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

type ParramLogin struct {
	Mobile      string `gorm:"size:255;not null;unique" json:"mobile"`
	DeviceID    string `gorm:"size:255;not null" json:"device_id"`
	DeviceToken string `gorm:"size:255;not null" json:"device_token"`
	Password    string `gorm:"size:100;not null;" json:"password"`
}

type ResultLogin struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	resultFail := ResultFail{}

	if err != nil {
		resultFail.Status = http.StatusUnprocessableEntity
		resultFail.Msg = err.Error()

		responses.JSON(w, http.StatusOK, resultFail)
		return
	}

	mdlogin := ParramLogin{}

	err = json.Unmarshal(body, &mdlogin)
	if err != nil {
		resultFail.Status = 1
		resultFail.Msg = err.Error()
		responses.JSON(w, http.StatusOK, resultFail)
		return
	}
	user := models.User{}
	user.Mobile = mdlogin.Mobile
	user.Password = mdlogin.Password

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		err = formaterror.FormatErrorLogin(err.Error())

		rs := formaterror.ReturnErr(err)

		responses.JSON(w, http.StatusOK, rs)
		return
	}
	token, err := server.SignIn(user.Mobile, user.Password)
	if err != nil {
		err = formaterror.FormatErrorLogin(err.Error())
		rs := formaterror.ReturnErr(err)

		responses.JSON(w, http.StatusOK, rs)
		return
	}

	//insert device token to database
	uid, err := auth.GetUIDByToken(token)

	loguserdevice := models.LogUserDevice{}
	_, err = loguserdevice.FindLogUserDeviceID(server.DB, mdlogin.DeviceID)
	if err != nil {
		loguserdevice.UserID = uid
		loguserdevice.LogType = 1
		loguserdevice.DeviceID = mdlogin.DeviceID
		loguserdevice.DeviceToken = mdlogin.DeviceToken
		loguserdevice.DeviceToken = mdlogin.DeviceToken
		loguserdevice.Prepare()
		err = loguserdevice.CreateLogUserDevice(server.DB)
		if err != nil {
			err := formaterror.ReturnErr(err)
			responses.JSON(w, http.StatusOK, err)
			return
		}

	}

	result := ResultLogin{}
	result.Status = http.StatusOK
	result.Data = token

	responses.JSON(w, http.StatusOK, result)
}

func (server *Server) SignIn(mobile, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("mobile = ?", mobile).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	return auth.CreateToken(user.ID)
}

func (server *Server) Logout(w http.ResponseWriter, r *http.Request) {
	//resultFail := ResultFail{}
	body, err := ioutil.ReadAll(r.Body)
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnSuccess()
		responses.JSON(w, http.StatusOK, err)
		return
	}

	login := ParramLogin{}
	err = json.Unmarshal(body, &login)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	logDevice := models.LogUserDevice{}

	_, err = logDevice.DeleteLogDevice(server.DB, login.DeviceID, uid)
	if err != nil {
		err := formaterror.ReturnErr(err)
		err.Msg = "DEVICE_NOT_FOUND"
		responses.JSON(w, http.StatusOK, err)
		return
	}

	result := formaterror.ReturnSuccess()
	responses.JSON(w, http.StatusOK, result)
}

func (server *Server) ForgotPass(w http.ResponseWriter, r *http.Request) {


	resultFail := ResultFail{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resultFail.Status = http.StatusUnprocessableEntity
		resultFail.Msg = err.Error()
		responses.JSON(w, http.StatusUnprocessableEntity, resultFail)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		resultFail.Status = http.StatusUnprocessableEntity
		resultFail.Msg = err.Error()
		responses.JSON(w, http.StatusUnprocessableEntity, resultFail)
		return
	}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", user.Email).Take(&user).Error
	if err != nil {
		resultFail.Status = http.StatusNoContent
		resultFail.Msg = err.Error()
		responses.JSON(w, http.StatusOK, resultFail)
		return
	}
	pass := "abc123"
	user.Password = pass
	uid := user.ID
	updatedUser, err := user.UpdatePasswordUser(server.DB, uint32(uid))
	if err != nil || updatedUser == nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

    contentMail := "Mật khẩu trên ứng dụng của bạn đã thay đổi thành " + pass + "vui lòng đăng nhập và đổi lại mật khẩu."
	send("Reset mật khẩu ứng dụng tâm linh",contentMail,user.Email,user.Password)
	result := ResultLogin{}
	result.Status = http.StatusOK
	result.Data = "success"
	responses.JSON(w, http.StatusOK, result)
}

func send(title string,body string, email string, password string) {
	from := "monopowernguyen@gmail.com"
	pass := "Toilaai@86"
	to := email

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: "+ title +"\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, visit http://foobarbazz.mailinator.com")
}

type ChangePass struct {
	OldPassword string `json:"old_pass"`
	NewPassword string `json:"new_pass"`
}

func (server *Server) ChangePassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	resultFail := ResultFail{}

	if err != nil {
		resultFail.Status = http.StatusUnprocessableEntity
		resultFail.Msg = err.Error()

		responses.JSON(w, http.StatusUnprocessableEntity, resultFail)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}

	//Check login

	changP := ChangePass{}
	err = json.Unmarshal(body, &changP)
	if err != nil {
		resultFail.Status = 1
		resultFail.Msg = err.Error()
		responses.JSON(w, http.StatusOK, resultFail)
		return
	}

	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		resultFail.Status = http.StatusUnprocessableEntity
		resultFail.Msg = err.Error()
		responses.JSON(w, http.StatusOK, resultFail)
		return
	}
	//Check old pass
	err = models.VerifyPassword(user.Password, changP.OldPassword)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		resultFail.Status = http.StatusUnprocessableEntity
		resultFail.Msg = "OLD_PASSWORD_WRONG"
		responses.JSON(w, http.StatusUnprocessableEntity, resultFail)
		return
	}
	user.Password = changP.NewPassword

	updatedUser, err := user.UpdatePasswordUser(server.DB, uint32(uid))
	if err != nil || updatedUser == nil {

		resultFail.Status = http.StatusUnprocessableEntity
		resultFail.Msg = err.Error()
		responses.JSON(w, http.StatusUnprocessableEntity, resultFail)
		return
	}
	result := ResultLogin{}
	result.Status = http.StatusOK
	result.Data = "success"
	responses.JSON(w, http.StatusOK, result)

}
