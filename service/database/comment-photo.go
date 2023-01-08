package database

func (db *appdbimpl) CommentPhoto(userID uint64, photoID uint64, comment string) error {
	_, err := db.c.Exec("INSERT INTO Comments (ownerid, photoid, comment) VALUES (?, ?, ?)", userID, photoID, comment)
	if err != nil {
		return err
	}

	return nil
}
