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

func (rt *_router) UnlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authID := utils.GetAuthorizationID(r.Header.Get("Authorization"))
	if authID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exists, err := rt.db.IdExists(authID)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlike-photo: error while checking if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseUint(ps.ByName("UserID"), 10, 64)

	if err != nil {
		ctx.Logger.WithError(err).Error("unlike-photo: userid parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if authID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	photoID, err := strconv.ParseUint(ps.ByName("PhotoID"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlike-photo: photoid parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	photoInfo, err := rt.db.GetPhotoInfo(photoID)

	if err != nil {
		if err == database.ErrPhotoNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("unlike-photo: error while getting the photo info from the database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	banned, err := rt.db.IsBanned(authID, photoInfo.OwnerID)

	if err != nil {
		ctx.Logger.WithError(err).Error("unlike-photo: error while checking if the requester is banned by the owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if banned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = rt.db.UnlikePhoto(userID, photoID)

	if err != nil {
		if errors.Is(err, database.ErrNotLiked) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("unlike-photo: error while removing the like")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
