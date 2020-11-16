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

func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post.Prepare()
	err = post.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)
}

func (server *Server) GetPosts(w http.ResponseWriter, r *http.Request) {

	post := models.Post{}
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

	items, err := post.FindAllPosts(server.DB, pid, p)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.IntroImage, "tamlinh")
			linkAvatar := server.getLink(v.Author.Avartar, "avatar")
			//fmt.Printf("dfdfdfd %v", linkImg)
			items[i].IntroImage = linkImg
			fmt.Printf("dfdfdfd %v", v)
			items[i].Author.Avartar = linkAvatar

		}
	}

	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

	//	result := formatresult.ReturnPost(posts)
	//	responses.JSON(w, http.StatusOK, result)
}

func (server *Server) GetSuggestPosts(w http.ResponseWriter, r *http.Request) {

	post := models.Post{}
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}

	p_id, errP := strconv.ParseUint(vars["post_id"], 10, 64)

	if errP != nil {
		err := formaterror.ReturnErr(errP)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	p, errP := strconv.ParseUint(vars["p"], 10, 64)

	if errP != nil {
		err := formaterror.ReturnErr(errP)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}

	items, err := post.FindSuggestPosts(server.DB, pid, p_id, p)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.IntroImage, "tamlinh")

			items[i].IntroImage = linkImg
			fmt.Printf("dfdfdfd %v", v)

		}
	}

	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

	//result := formatresult.ReturnPost(posts)
	//responses.JSON(w, http.StatusOK, result)
}

func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	post := models.Post{}

	item, err := post.FindPostByID(server.DB, pid)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	linkImg := server.getLink(item.IntroImage, "tamlinh")
	//fmt.Printf("dfdfdfd %v", linkImg)
	item.IntroImage = linkImg
	linkAvatar := server.getLink(item.Author.Avartar, "avatar")
	item.Author.Avartar = linkAvatar

	var rs interface{}
	rs = item

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

}

func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {

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
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	// If a user attempt to update a post not belonging to him
	if uid != post.AuthorID {
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
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != postUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	postUpdate.Prepare()
	err = postUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate.ID = post.ID //this is important to tell the model the post id to update, the other update field are set above

	postUpdated, err := postUpdate.UpdateAPost(server.DB)

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated)
}

func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {

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
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = post.DeleteAPost(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
