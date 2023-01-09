package database

func (db *appdbimpl) DeleteComment(commentID uint64) error {
	sqlres, err := db.c.Exec("DELETE FROM Comments WHERE id = ?", commentID)
	if err != nil {
		return err
	}
	rowsAffected, err := sqlres.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCommentNotFound
	}

	return nil
}
