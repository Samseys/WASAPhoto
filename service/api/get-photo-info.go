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

func (rt *_router) GetPhotoInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authID := utils.GetAuthorizationID(r.Header.Get("Authorization"))
	if authID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exists, err := rt.db.IdExists(authID)
	if err != nil {
		ctx.Logger.WithError(err).Error("get-photo-info: error while checking if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photoID, err := strconv.ParseUint(ps.ByName("PhotoID"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("get-photo-info: parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	photoInfo, err := rt.db.GetPhotoInfoForFrontend(photoID)

	if err != nil {
		if errors.Is(err, database.ErrPhotoNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("get-photo-info: error while getting the photo info from the database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	banned, err := rt.db.IsBanned(authID, photoInfo.Owner.ID)

	if err != nil {
		ctx.Logger.WithError(err).Error("get-photo-info: error while checking if the requester is banned by the owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if banned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(photoInfo)
}
