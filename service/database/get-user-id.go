package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetUserID(username string) (uint64, error) {
	var id uint64
	err := db.c.QueryRow("SELECT id FROM Users WHERE name=?", username).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrUserProfileNotFound
	}
	return id, err
}
