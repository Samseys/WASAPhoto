package api

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
	"me.samsey/wasa-photos/service/utils"
)

func (rt *_router) DeletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authID := utils.GetAuthorizationID(r.Header.Get("Authorization"))
	if authID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exists, err := rt.db.IdExists(authID)
	if err != nil {
		ctx.Logger.WithError(err).Error("delete-photo: error while checking if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photoID, err := strconv.ParseUint(ps.ByName("PhotoID"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("delete-photo: photoid parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	photoInfo, err := rt.db.GetPhotoInfo(photoID)

	if err != nil {
		if errors.Is(err, database.ErrPhotoNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("delete-photo: error while getting the photo info from the database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if authID != photoInfo.OwnerID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = rt.db.DeletePhoto(photoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("delete-photo: error deleting the photo from the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	photoPath := filepath.Join(database.PhotoPath, strconv.FormatUint(authID, 10), strconv.FormatUint(photoID, 10)+photoInfo.Extension)
	err = os.Remove(photoPath)
	if err != nil {
		ctx.Logger.WithError(err).Error("delete-photo: error deleting the photo from the filesystem")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
