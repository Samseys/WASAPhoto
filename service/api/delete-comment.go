package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
	"me.samsey/wasa-photos/service/utils"
)

func (rt *_router) DeleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authID := utils.GetAuthorizationID(r.Header.Get("Authorization"))
	if authID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exists, err := rt.db.IdExists(authID)
	if err != nil {
		ctx.Logger.WithError(err).Error("delete-comment: error while checking if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photoID, err := strconv.ParseUint(ps.ByName("PhotoID"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("delete-comment: photoid parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	photoInfo, err := rt.db.GetPhotoInfo(photoID)

	if err != nil {
		if err == database.ErrPhotoNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("delete-photo: error while getting the photo info from the database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	banned, err := rt.db.IsBanned(authID, photoInfo.OwnerID)

	if err != nil {
		ctx.Logger.WithError(err).Error("delete-comment: error while checking if the requester is banned by the owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if banned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	commentID, err := strconv.ParseUint(ps.ByName("CommentID"), 10, 64)

	if err != nil {
		ctx.Logger.WithError(err).Error("delete-comment: parsing the comment id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comment, err := rt.db.GetComment(commentID)

	if err != nil {
		if errors.Is(database.ErrCommentNotFound, err) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("delete-comment: error getting the comment info")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if authID != comment.Owner.ID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = rt.db.DeleteComment(commentID)

	if err != nil {
		if errors.Is(err, database.ErrNotFollowed) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("delete-comment: error while removing the ban from the database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}