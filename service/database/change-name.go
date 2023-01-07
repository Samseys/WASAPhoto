package database

import "strings"

func (db *appdbimpl) ChangeName(username string, identifier int) error {
	_, err := db.c.Exec("UPDATE Users SET name = ? WHERE id = ?", username, identifier)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return ErrUsernameAlradyTaken
		} else {
			return err
		}
	}

	return nil
}
