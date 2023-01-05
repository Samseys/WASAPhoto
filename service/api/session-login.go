package api

import (
	"net/http"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
)

type Name struct {
	Name string `json:"name"`
}

type Id struct {
	Identifier int `json:"identifier"`
}

func (rt *_router) Session(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var name Name
	err := json.NewDecoder(r.Body).Decode(&name)
	if err != nil {
		// The body was not a parseable JSON, reject it
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if name.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	identifier, err := rt.db.Login(name.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error logging user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var response = Id{Identifier: identifier}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
