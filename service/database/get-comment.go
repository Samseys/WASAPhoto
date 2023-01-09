package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetComment(commentID uint64) (Comment, error) {
	var comment Comment
	err := db.c.QueryRow("SELECT id, ownerid, creationdate, comment FROM Comments WHERE id = ?", commentID).Scan(&comment.ID, &comment.Owner.ID, &comment.CreationDate, &comment.Comment)
	if errors.Is(err, sql.ErrNoRows) {
		return comment, ErrCommentNotFound
	}
	err = db.c.QueryRow("SELECT name FROM Users WHERE id = ?", comment.Owner.ID).Scan(&comment.Owner.Username)
	return comment, err
}
