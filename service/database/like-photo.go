package database

func (db *appdbimpl) LikePhoto(userID uint64, photoID uint64) error {
	sqlres, err := db.c.Exec("INSERT INTO Likes (photoid, userid) VALUES (?, ?) ON CONFLICT DO NOTHING", photoID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := sqlres.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrAlreadyLiked
	}

	return nil
}
