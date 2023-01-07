package database

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (db *appdbimpl) IdExistsAndCompare(r *http.Request, ps httprouter.Params) (bool, int) {
	splitted := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(splitted) != 2 {
		return false, -1
	}
	id, err := strconv.Atoi(splitted[1])
	if err != nil {
		return false, -1
	}
	var exists bool
	db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM Users WHERE id = ?)", id).Scan(&exists)
	if ps == nil && exists {
		return true, id
	}

	useridFromPath, err := strconv.Atoi(ps.ByName("UserID"))
	if err == nil && exists && id == useridFromPath {
		return true, id
	} else {
		return false, -1
	}
}

// return False if the requester is not logged or not banned by the other user, True otherwise
func (db *appdbimpl) IsBanned(r *http.Request, otherUserID int) bool {
	splitted := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(splitted) != 2 {
		return false
	}

	id, err := strconv.Atoi(splitted[1])
	if err != nil {
		return false
	}
	var exists bool
	db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM Bans WHERE userid = ? AND bannedid = ?)", otherUserID, id).Scan(&exists)
	return exists
}
