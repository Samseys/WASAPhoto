package database

import "errors"

var ErrNotFollowed = errors.New("the user is not followed")

func (db *appdbimpl) Unfollow(followerID int, followedID int) error {
	sqlres, err := db.c.Exec("DELETE FROM Follows WHERE followerid = ? AND followedid = ?", followerID, followedID)
	if err != nil {
		return err
	}
	rowsAffected, _ := sqlres.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFollowed
	}

	return nil
}
