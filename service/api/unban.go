package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
)

func (rt *_router) UnbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	existsAndEqual, id := rt.db.IdExistsAndCompare(r, ps)
	if !existsAndEqual {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	otherUserID, err := strconv.ParseUint(ps.ByName("OtherUserID"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("unban: parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := rt.db.IdExists(otherUserID)

	if err != nil {
		ctx.Logger.WithError(err).Error("unban: error while checking if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rt.db.Unban(id, otherUserID)

	if err != nil {
		if errors.Is(err, database.ErrNotFollowed) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			ctx.Logger.WithError(err).Error("unban: error while removing the ban from the database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
