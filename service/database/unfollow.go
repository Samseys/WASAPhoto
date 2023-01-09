package database

func (db *appdbimpl) UnfollowUser(followerID uint64, followedID uint64) error {
	sqlres, err := db.c.Exec("DELETE FROM Follows WHERE followerid = ? AND followedid = ?", followerID, followedID)
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
