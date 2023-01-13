package api

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
	"me.samsey/wasa-photos/service/utils"
)

func (rt *_router) GetPhotoFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authID := utils.GetAuthorizationID(r.Header.Get("Authorization"))
	if authID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exists, err := rt.db.IdExists(authID)
	if err != nil {
		ctx.Logger.WithError(err).Error("get-photo-file: error while checking if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photoID, err := strconv.ParseUint(ps.ByName("PhotoID"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("get-photo-file: parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	photoInfo, err := rt.db.GetPhotoInfo(photoID)

	if err != nil {
		if errors.Is(err, database.ErrPhotoNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("get-photo-file: error while getting the photo info from the database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	banned, err := rt.db.IsBanned(authID, photoInfo.OwnerID)

	if err != nil {
		ctx.Logger.WithError(err).Error("get-photo-file: error while checking if the requester is banned by the owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if banned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	photoPath := filepath.Join(database.PhotoPath, strconv.FormatUint(photoInfo.OwnerID, 10), strconv.FormatUint(photoID, 10)+photoInfo.Extension)
	photo, err := os.ReadFile(photoPath)
	if err != nil {
		ctx.Logger.WithError(err).Error("get-photo-file: error while getting the file from the disk")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "image/"+strings.ReplaceAll(photoInfo.Extension, ".", ""))
	_, err = w.Write(photo)
	if err != nil {
		ctx.Logger.WithError(err).Error("get-photo-file: error while writing the image into the response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
