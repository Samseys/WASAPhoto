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

	otherUserID, err := strconv.Atoi(ps.ByName("OtherUserID"))
	if err != nil {
		ctx.Logger.Error("follow: parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exists := rt.db.IdExists(otherUserID)

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
			ctx.Logger.WithError(err).Error("can't process the follow request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
