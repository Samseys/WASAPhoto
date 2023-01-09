package database

func (db *appdbimpl) GetCommentOwner(commentID uint64) (uint64, error) {
	var ownerID uint64
	err := db.c.QueryRow("SELECT ownerid FROM Comments WHERE id = ?", commentID).Scan(&ownerID)

	return ownerID, err
}
