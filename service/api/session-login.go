package api

import (
	"net/http"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
)

func (rt *_router) Session(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var name database.Name
	err := json.NewDecoder(r.Body).Decode(&name)
	defer r.Body.Close()
	if err != nil {
		ctx.Logger.WithError(err).Error("session: error decoding user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if name.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	identifier, err := rt.db.Login(name.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("session: error logging user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var response = database.UserId{Identifier: identifier}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
