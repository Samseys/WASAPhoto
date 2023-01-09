package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
	"me.samsey/wasa-photos/service/utils"
)

func (rt *_router) BanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authID := utils.GetAuthorizationID(r.Header.Get("Authorization"))
	if authID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exists, err := rt.db.IdExists(authID)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban: error while checking if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseUint(ps.ByName("UserID"), 10, 64)

	if err != nil {
		ctx.Logger.WithError(err).Error("ban: userid parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if authID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	otherUserID, err := strconv.ParseUint(ps.ByName("OtherUserID"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban: otheruserid parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userID == otherUserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	exists, err = rt.db.IdExists(otherUserID)

	if err != nil {
		ctx.Logger.WithError(err).Error("ban: error while checking if the other user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rt.db.BanUser(userID, otherUserID)

	if err != nil {
		if errors.Is(err, database.ErrAlreadyBanned) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			ctx.Logger.WithError(err).Error("ban: errore while inserting the ban")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
