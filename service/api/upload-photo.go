package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
	"me.samsey/wasa-photos/service/utils"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authID := utils.GetAuthorizationID(r.Header.Get("Authorization"))
	if authID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exists, err := rt.db.IdExists(authID)
	if err != nil {
		ctx.Logger.WithError(err).Error("follow: error while checking if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: error parsing multipart header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mainComment := r.FormValue("MainComment")

	if mainComment == "" {
		ctx.Logger.WithError(err).Error("upload-photo: empty comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, fileheader, err := r.FormFile("UploadedPhoto")
	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: error while reading the photo")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(fileheader.Filename)

	photoid, err := rt.db.UploadPhoto(authID, mainComment, ext)

	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: database error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: error while reading the file extension")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: error while resetting the pointer")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	photoPath := filepath.Join(database.PhotoPath, strconv.FormatUint(authID, 10))
	err = os.MkdirAll(photoPath, os.ModePerm)
	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: error while creating the folder path")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(filepath.Join(photoPath, strconv.FormatUint(photoid, 10)+ext))
	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: error while creating the file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		ctx.Logger.WithError(err).Error("upload-photo: error while copying the file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := database.PhotoID{ID: photoid}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
