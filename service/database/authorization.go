package database

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (db *appdbimpl) IdExistsAndCompare(r *http.Request, ps httprouter.Params) (bool, uint64) {
	splitted := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(splitted) != 2 {
		return false, 0
	}
	id, err := strconv.Atoi(splitted[1])
	if err != nil {
		return false, 0
	}
	var exists bool
	err = db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM Users WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false, 0
	}

	if ps == nil && exists {
		return true, uint64(id)
	}

	useridFromPath, err := strconv.Atoi(ps.ByName("UserID"))
	if err == nil && exists && id == useridFromPath {
		return true, uint64(id)
	} else {
		return false, 0
	}
}

func (db *appdbimpl) IdExists(id uint64) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM Users WHERE id = ?)", id).Scan(&exists)

	return exists, err
}

// return False if the requester is not logged or not banned by the other user, True otherwise
func (db *appdbimpl) IsBanned(r *http.Request, otherUserID uint64) (bool, error) {
	splitted := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(splitted) != 2 {
		return false, errors.New("format not valid")
	}

	id, err := strconv.Atoi(splitted[1])
	if err != nil {
		return false, err
	}
	var exists bool
	err = db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM Bans WHERE userid = ? AND bannedid = ?)", otherUserID, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
