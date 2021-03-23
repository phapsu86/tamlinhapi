package controllers

import (
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/phapsu86/tamlinh/api/auth"
	"github.com/phapsu86/tamlinh/api/models"
	"github.com/phapsu86/tamlinh/api/responses"
	"github.com/phapsu86/tamlinh/api/utils/formaterror"
	"github.com/phapsu86/tamlinh/api/utils/formatresult"
)

func (server *Server) CreatePostComment(w http.ResponseWriter, r *http.Request) {

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
	post := models.PostComment{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusUnprocessableEntity, err)

		return
	}

	fmt.Printf("%v\n", post)
	post.Prepare()
	err = post.Validate()
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	post.UserID = uid

	_, err = post.SavePostComment(server.DB)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}

	result := formaterror.ReturnSuccess()
	responses.JSON(w, http.StatusCreated, result)
}

func (server *Server) GetPostCommentByID(w http.ResponseWriter, r *http.Request) {

	post := models.PostComment{}
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	p, errP := strconv.ParseUint(vars["p"], 10, 64)

	if errP != nil {
		err := formaterror.ReturnErr(errP)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}

	posts, err := post.FindAllPostComment(server.DB, pid, p)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	if len(posts) > 0 {
		for i, v := range posts {
			linkImg := server.getLink(v.User.Avartar, "avatar")
			//fmt.Printf("dfdfdfd %v", linkImg)
			posts[i].User.Avartar = linkImg
			fmt.Printf("dfdfdfd %v", v)

		}
	}

	var rs interface{}
	rs = posts
	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)
	// result := formatresult.ReturnPostComment(posts)
	// responses.JSON(w, http.StatusOK, result)
}
