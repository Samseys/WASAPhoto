package api

import (
	"errors"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
	"me.samsey/wasa-photos/service/utils"
)

func (rt *_router) ChangeName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authid := utils.GetAuthorizationID(r.Header.Get("Authorization"))
	if authid == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exists, err := rt.db.IdExists(authid)
	if err != nil {
		ctx.Logger.WithError(err).Error("change-name: error while checking if the user exists")
		w.WriteHeader(http.StatusBadRequest)
	}

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userid, err := strconv.ParseUint(ps.ByName("UserID"), 10, 64)

	if err != nil {
		ctx.Logger.WithError(err).Error("change-name: userid parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if authid != userid {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var name database.Name
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&name)
	if err != nil {
		ctx.Logger.WithError(err).Error("change-name: error decoding JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if name.Name == "" {
		ctx.Logger.Error("change-name: name empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name.Name = utils.MakeAlphaNumeric(name.Name)

	err = rt.db.ChangeName(name.Name, userid)

	if err != nil {
		if errors.Is(err, database.ErrUsernameAlreadyTaken) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			ctx.Logger.WithError(err).Error("change-name: error while inserting the name change in the database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
