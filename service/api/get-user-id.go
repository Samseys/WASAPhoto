package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
)

func (rt *_router) GetUserID(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	name := ps.ByName("Username")
	id, err := rt.db.GetUserID(name)
	if err != nil {
		ctx.Logger.WithError(err).Error("getuserid: error getting user id from db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := database.UserID{ID: id}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
