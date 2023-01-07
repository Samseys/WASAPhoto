package database

import "strings"

func (db *appdbimpl) ChangeName(username string, id int) error {
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
