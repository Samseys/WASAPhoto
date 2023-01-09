package database

func (db *appdbimpl) UnlikePhoto(userID uint64, photoID uint64) error {
	sqlres, err := db.c.Exec("DELETE FROM Likes WHERE userid = ? AND photoid = ?", userID, photoID)
	if err != nil {
		return err
	}
	rowsAffected, err := sqlres.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotLiked
	}

	return nil
}
