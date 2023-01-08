package database

import "errors"

var ErrNotLiked = errors.New("the photo is not liked")

func (db *appdbimpl) UnlikePhoto(userID uint64, photoID uint64) error {
	sqlres, err := db.c.Exec("DELETE FROM Like WHERE userid = ? AND photoid = ?", userID, photoID)
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
