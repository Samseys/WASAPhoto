package database

import "errors"

var ErrNotBanned = errors.New("the user is not banned")

func (db *appdbimpl) Unban(userID int, bannedID int) error {
	sqlres, err := db.c.Exec("DELETE FROM Bans WHERE userid = ? AND bannedid = ?", userID, bannedID)
	if err != nil {
		return err
	}
	rowsAffected, _ := sqlres.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFollowed
	}

	return nil
}
