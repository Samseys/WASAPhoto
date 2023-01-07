package api

import (
	"errors"
	"net/http"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
)

func (rt *_router) ChangeName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	existsAndEqual, id := rt.db.IdExistsAndCompare(r, ps)
	if !existsAndEqual {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var name database.Name
	err := json.NewDecoder(r.Body).Decode(&name)
	defer r.Body.Close()
	if err != nil {
		ctx.Logger.WithError(err).Error("change name: error decoding JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if name.Name == "" {
		ctx.Logger.Error("change name: error validating JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.ChangeName(name.Name, id)

	if err != nil {
		if errors.Is(err, database.ErrUsernameAlradyTaken) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			ctx.Logger.WithError(err).Error("can't process the change name request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
