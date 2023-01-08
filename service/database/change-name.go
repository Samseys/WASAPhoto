package database

import (
	"errors"
	"strings"
)

var ErrUsernameAlradyTaken = errors.New("another user already has this username")

func (db *appdbimpl) ChangeName(username string, id uint64) error {
	_, err := db.c.Exec("UPDATE Users SET name = ? WHERE id = ?", username, id)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return ErrUsernameAlradyTaken
		} else {
			return err
		}
	}

	return nil
}
