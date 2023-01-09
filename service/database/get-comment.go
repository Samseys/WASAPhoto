package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetComment(commentID uint64) (Comment, error) {
	var comment Comment
	err := db.c.QueryRow("SELECT Comments.id as commentids, Users.id as ownerid, Users.name, Comments.creationdate, Comments.comment FROM Comments INNER JOIN Users ON Comments.ownerid = Users.id WHERE Comments.id = ?", commentID).Scan(&comment.ID, &comment.Owner.ID, &comment.Owner.Username, &comment.CreationDate, &comment.Comment)
	if errors.Is(err, sql.ErrNoRows) {
		return comment, ErrCommentNotFound
	}
	return comment, err
}
