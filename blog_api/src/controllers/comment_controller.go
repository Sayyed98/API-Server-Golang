package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"blog_api/src/auth"
	"blog_api/src/models"
	"blog_api/src/responses"
	"blog_api/src/utils/formaterror"

	"github.com/gorilla/mux"
)

func (server *Server) CreateComment(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	comment := models.Comment{}
	err = json.Unmarshal(body, &comment)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	commentCreated, err := comment.SaveComment(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, commentCreated.ID))
	responses.JSON(w, http.StatusCreated, commentCreated)
}

func (server *Server) GetComments(w http.ResponseWriter, r *http.Request) {

	comments := models.Comment{}

	comment, err := comments.CommentAll(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, comment)
}

func (server *Server) GetComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	comment := models.Comment{}

	commentReceived, err := comment.CommentByID(server.DB, cid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, commentReceived)
}

func (server *Server) UpdateComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	comment := models.Comment{}
	err = server.DB.Debug().Model(models.Comment{}).Where("id = ?", cid).Take(&comment).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	commentUpdate := models.Comment{}
	err = json.Unmarshal(body, &commentUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	commentUpdated, err := commentUpdate.UpdateComment(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, commentUpdated)
}

func (server *Server) DeleteComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	comment := models.Comment{}

	// Is the authenticated user, the owner of this post?
	if uid != comment.Author.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = comment.DeleteComment(server.DB, cid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", cid))
	responses.JSON(w, http.StatusNoContent, "")
}
