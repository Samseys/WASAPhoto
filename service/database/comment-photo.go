package database

func (db *appdbimpl) CommentPhoto(userID uint64, photoID uint64, comment string) (uint64, error) {
	sqlres, err := db.c.Exec("INSERT INTO Comments (ownerid, photoid, comment) VALUES (?, ?, ?)", userID, photoID, comment)
	if err != nil {
		return 0, err
	}
	commentID, err := sqlres.LastInsertId()
	return uint64(commentID), err
}
