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
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	existsAndEqual, id := rt.db.IdExistsAndCompare(r, nil)
	if !existsAndEqual {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mainComment := r.FormValue("mainComment")

	if mainComment == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, fileheader, err := r.FormFile("uploadedImage")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(fileheader.Filename)

	photoid, err := rt.db.UploadImage(id, mainComment, ext)

	if err != nil {
		ctx.Logger.WithError(err).Error("can't process the image upload request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't process the image upload request")
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	imagePath := filepath.Join(database.ImagePath, strconv.FormatUint(id, 10))
	err = os.MkdirAll(imagePath, os.ModePerm)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't process the image upload request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(filepath.Join(imagePath, strconv.FormatUint(photoid, 10)+ext))
	if err != nil {
		ctx.Logger.WithError(err).Error("can't process the image upload request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't process the image upload request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := database.PhotoID{ID: photoid}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
