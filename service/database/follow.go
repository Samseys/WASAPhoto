package database

import (
	"errors"
)

var ErrAlreadyFollowed = errors.New("the user is already followed")

func (db *appdbimpl) Follow(followerID int, followedID int) error {
	sqlres, err := db.c.Exec("INSERT INTO Follows (followerid, followedid) VALUES (?, ?) ON CONFLICT DO NOTHING", followerID, followedID)
	if err != nil {
		return err
	}
	rowsAffected, _ := sqlres.RowsAffected()
	if rowsAffected == 0 {
		return ErrAlreadyFollowed
	}

	return nil
}
