package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
	"me.samsey/wasa-photos/service/utils"
)

func (rt *_router) CommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authID := utils.GetAuthorizationID(r.Header.Get("Authorization"))
	if authID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exists, err := rt.db.IdExists(authID)
	if err != nil {
		ctx.Logger.WithError(err).Error("comment-photo: error while checking if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photoID, err := strconv.ParseUint(ps.ByName("PhotoID"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("comment-photo: photoid parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	photoInfo, err := rt.db.GetPhotoInfo(photoID)

	if err != nil {
		if errors.Is(err, database.ErrPhotoNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("comment-photo: error while getting the photo info from the database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	banned, err := rt.db.IsBanned(authID, photoInfo.OwnerID)

	if err != nil {
		ctx.Logger.WithError(err).Error("comment-photo: error while checking if the requester is banned by the photo owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if banned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var comment database.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	defer r.Body.Close()
	if err != nil {
		ctx.Logger.WithError(err).Error("comment-photo: error decoding the JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if comment.Comment == "" {
		ctx.Logger.Error("comment-photo: empty comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	commentID, err := rt.db.CommentPhoto(authID, photoID, comment.Comment)

	if err != nil {
		ctx.Logger.WithError(err).Error("comment-photo: error while inserting the comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := rt.db.GetComment(commentID)
	if err != nil {
		ctx.Logger.WithError(err).Error("comment-photo: error while getting the new comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
