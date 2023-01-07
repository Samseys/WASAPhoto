package database

import (
	"errors"
)

var ErrAlreadyBanned = errors.New("the user is already banned")

func (db *appdbimpl) Ban(userID int, bannedID int) error {
	sqlres, err := db.c.Exec("INSERT INTO Bans (userid, bannedid) VALUES (?, ?) ON CONFLICT DO NOTHING", userID, bannedID)
	if err != nil {
		return err
	}

	rowsAffected, _ := sqlres.RowsAffected()
	if rowsAffected == 0 {
		return ErrAlreadyBanned
	}

	return nil
}
