package database

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (db *appdbimpl) IdExistsAndEqual(r *http.Request, ps httprouter.Params) (bool, int) {
	splitted := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(splitted) != 2 {
		return false, -1
	}

	var id, _ = strconv.Atoi(splitted[1])
	var exists bool
	db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM Users WHERE id = ?)", id).Scan(&exists)
	useridFromPath, _ := strconv.Atoi(ps.ByName("UserID"))

	if exists && id == useridFromPath {
		return true, id
	} else {
		return false, -1
	}
}
