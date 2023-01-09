package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/utils"
)

func (rt *_router) Stream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authID := utils.GetAuthorizationID(r.Header.Get("Authorization"))
	if authID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exists, err := rt.db.IdExists(authID)
	if err != nil {
		ctx.Logger.WithError(err).Error("stream: error while checking if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	stream, err := rt.db.GetStream(authID)
	if err != nil {
		ctx.Logger.WithError(err).Error("stream: error while getting the stream data from the db")
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(stream)
}
