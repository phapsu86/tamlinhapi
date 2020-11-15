package formaterror

import (
	"errors"
	"strings"
)

type ErrMsg struct { 
	Status   int  `json:"status"` 
	Msg string  `json:"msg"` 
	
}

type ErrMsgUpload struct { 
	Status   int  `json:"status"` 
	Msg string  `json:"msg"`
	Url string  `json:"url"`
	
}


func ReturnErr(err error) (ErrMsg) {
	errMs := ErrMsg{}
	errMs.Status = 1
	errMs.Msg = err.Error()
	return errMs
}

func ReturnSuccess() (ErrMsg){
	errMs := ErrMsg{}
	errMs.Status = 200
	errMs.Msg ="success"
	return errMs

}

func ReturnUploadSuccess() (ErrMsgUpload){
	errMs := ErrMsgUpload{}
	errMs.Status = 200
	errMs.Msg ="success"

	return errMs

}

func FormatErrorLogin(err string) error {

	if strings.Contains(err, "nickname") {
		return errors.New("Nickname Already Taken")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}

	if strings.Contains(err, "title") {
		return errors.New("Title Already Taken")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}
	return errors.New("Incorrect Details")
}

