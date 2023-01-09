package database

func (db *appdbimpl) CommentExists(commentID uint64) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM Comments WHERE id = ?)", commentID).Scan(&exists)
	return exists, err
}
