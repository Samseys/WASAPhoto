package database

import "errors"

func (db *appdbimpl) ChangeName(username string, identifier int) error {
	sqlres, _ := db.c.Exec("UPDATE Users SET name = ? WHERE id = ?", username, identifier)
	if sqlres == nil {
		return errors.New("another user already has this username")
	}

	return nil
}
