package database

func (db *appdbimpl) UploadImage(userID uint64, mainComment string, ext string) (uint64, error) {
	sqlres, err := db.c.Exec("INSERT INTO Photos (ownerid, maincomment, extension) VALUES (?, ?, ?)", userID, mainComment, ext)
	if err != nil {
		return 0, err
	}

	photoID, err := sqlres.LastInsertId()

	if err != nil {
		return 0, err
	}

	return uint64(photoID), nil
}
