package middlewares

import (
	//"errors"
	"net/http"

	"github.com/victorsteven/fullstack/api/auth"
	"github.com/victorsteven/fullstack/api/responses"
	"github.com/victorsteven/fullstack/api/utils/formaterror"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			err:= formaterror.ReturnErr(err)
			responses.JSON(w, http.StatusOK,err)
		//	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
