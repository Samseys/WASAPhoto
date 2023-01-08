package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
)

func (rt *_router) UnfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	existsAndEqual, id := rt.db.IdExistsAndCompare(r, ps)
	if !existsAndEqual {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	otherUserID, err := strconv.ParseUint(ps.ByName("OtherUserID"), 10, 64)
	if err != nil {
		ctx.Logger.Error("unfollow: parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := rt.db.IdExists(otherUserID)

	if err != nil {
		ctx.Logger.Error("can't process the unfollow request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rt.db.Unfollow(id, otherUserID)

	if err != nil {
		if errors.Is(err, database.ErrNotFollowed) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			ctx.Logger.WithError(err).Error("can't process the unfollow request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
