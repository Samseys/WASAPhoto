package database

import "errors"

var ErrNotBanned = errors.New("the user is not banned")

func (db *appdbimpl) UnbanUser(userID uint64, bannedID uint64) error {
	sqlres, err := db.c.Exec("DELETE FROM Bans WHERE userid = ? AND bannedid = ?", userID, bannedID)
	if err != nil {
		return err
	}
	rowsAffected, err := sqlres.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFollowed
	}

	return nil
}
