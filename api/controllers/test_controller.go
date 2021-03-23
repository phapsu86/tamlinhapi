package controllers

import (
	
	"net/http"


	
	"github.com/phapsu86/tamlinh/api/models"
	"github.com/phapsu86/tamlinh/api/responses"
	
)
//Chưc năng get test
func (server *Server) GetTest(w http.ResponseWriter, r *http.Request) {

	post := models.Hung{}

	posts, err := post.FindAllTest(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}
