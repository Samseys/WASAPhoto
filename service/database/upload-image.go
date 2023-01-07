package database

func (db *appdbimpl) UploadImage(userID int, mainComment string, ext string) (int, error) {
	sqlres, err := db.c.Exec("INSERT INTO Photos (ownerid, maincomment, extension) VALUES (?, ?, ?)", userID, mainComment, ext)
	if err != nil {
		return -1, err
	}

	photoID, err := sqlres.LastInsertId()

	if err != nil {
		return -1, err
	}

	return int(photoID), nil
}
