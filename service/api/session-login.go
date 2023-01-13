package api

import (
	"net/http"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
	"me.samsey/wasa-photos/service/utils"
)

func (rt *_router) Session(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var name database.Name
	err := json.NewDecoder(r.Body).Decode(&name)
	defer r.Body.Close()
	if err != nil {
		ctx.Logger.WithError(err).WithError(err).Error("session: error decoding user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !utils.CheckName(name.Name) {
		ctx.Logger.Error("session: input not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := rt.db.Login(name.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("session: error logging user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := database.UserID{ID: id}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
