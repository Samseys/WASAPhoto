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
		ctx.Logger.Error("comment-photo: parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err = rt.db.PhotoExists(photoID)

	if err != nil {
		ctx.Logger.Error("comment-photo: error while checking if the photo exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ownerID, err := rt.db.GetPhotoOwner(photoID)
	if err != nil {
		ctx.Logger.Error("comment-photo: error while getting the photo owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	banned, err := rt.db.IsBanned(authID, ownerID)

	if err != nil {
		ctx.Logger.Error("comment-photo: error while checking if the requester is banned by the owner")
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
		ctx.Logger.WithError(err).Error("comment-photo: error decoding ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if comment.Comment == "" {
		ctx.Logger.WithError(err).Error("comment-photo: empty comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.CommentPhoto(authID, photoID, comment.Comment)

	if err != nil {
		if errors.Is(err, database.ErrAlreadyLiked) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			ctx.Logger.Error("comment-photo: error while inserting the comment")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
