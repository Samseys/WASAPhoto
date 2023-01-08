package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"me.samsey/wasa-photos/service/api/reqcontext"
	"me.samsey/wasa-photos/service/database"
)

func (rt *_router) GetUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userid, err := strconv.ParseUint(ps.ByName("UserID"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("get user profile: parameter not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isBanned, err := rt.db.IsBanned(r, userid)
	if err != nil {
		ctx.Logger.WithError(err).Error("get user profile: error while checking if the requester is banned by the other user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isBanned {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	profile, err := rt.db.GetUserProfile(userid)

	if err != nil {
		if errors.Is(err, database.ErrUserProfileNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("get user profile: error while getting the user profile")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(profile)
}
