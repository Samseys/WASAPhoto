package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
)

func (rt *_router) LikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	existsAndEqual, id := rt.db.IdExistsAndCompare(r, ps)
	if !existsAndEqual {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photoID, err := strconv.ParseUint(ps.ByName("PhotoID"), 10, 64)
	if err != nil {
		ctx.Logger.Error("like-photo: parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := rt.db.PhotoExists(photoID)

	if err != nil {
		ctx.Logger.Error("like-photo: error while checking if the photo exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ownerID, err := rt.db.GetPhotoOwner(photoID)
	if err != nil {
		ctx.Logger.Error("like-photo: error while getting the photo owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	banned, err := rt.db.IsBanned(r, ownerID)

	if err != nil {
		ctx.Logger.Error("like-photo: error while checking if the requester is banned by the owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if banned {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = rt.db.LikePhoto(id, photoID)

	if err != nil {
		if errors.Is(err, database.ErrAlreadyLiked) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			ctx.Logger.Error("like-photo: error while inserting the like")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
