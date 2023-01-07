package database

import (
	"errors"
)

var ErrAlreadyFollowed = errors.New("aa")

func (db *appdbimpl) Follow(followeridentifier int, followedidentifier int) error {
	sqlres, err := db.c.Exec("INSERT INTO Follows (followerid, followedid) VALUES (?, ?) ON CONFLICT DO NOTHING", followeridentifier, followedidentifier)
	if err != nil {
		return err
	}
	var rowsAffected, _ = sqlres.RowsAffected()
	if rowsAffected == 0 {
		return ErrAlreadyFollowed
	}

	return nil
}
