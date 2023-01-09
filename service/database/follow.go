package database

func (db *appdbimpl) FollowUser(followerID uint64, followedID uint64) error {
	sqlres, err := db.c.Exec("INSERT INTO Follows (followerid, followedid) VALUES (?, ?) ON CONFLICT DO NOTHING", followerID, followedID)
	if err != nil {
		return err
	}
	rowsAffected, err := sqlres.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrAlreadyFollowed
	}

	return nil
}
