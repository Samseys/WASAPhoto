package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetPhotoInfo(photoID uint64) (Photo, error) {
	var photo Photo
	err := db.c.QueryRow("SELECT id, ownerid, extension, creationdate FROM Photos WHERE id = ?", photoID).Scan(&photo.ID, &photo.OwnerID, &photo.Extension, &photo.CreationDate)
	if errors.Is(err, sql.ErrNoRows) {
		return photo, ErrPhotoNotFound
	}

	return photo, nil
}

func (db *appdbimpl) GetPhotoInfoForFrontend(photoID uint64) (PhotoForFrontend, error) {
	var photo PhotoForFrontend
	err := db.c.QueryRow("SELECT id, ownerid, creationdate FROM Photos WHERE id = ?", photoID).Scan(&photo.ID, &photo.Owner.ID, &photo.CreationDate)
	if errors.Is(err, sql.ErrNoRows) {
		return photo, ErrPhotoNotFound
	}
	rows, err := db.c.Query("SELECT Comments.id, Users.id, Users.name, Comments.creationdate, Comments.comment FROM Comments INNER JOIN Users ON Comments.ownerid = Users.id WHERE photoid = ?", photoID)
	if err != nil {
		return photo, err
	}

	for rows.Next() {
		var comment Comment
		err = rows.Scan(comment.ID, comment.Owner.ID)
		if err != nil {
			return photo, err
		}
		photo.Comments = append(photo.Comments, comment)
	}

	rows.Close()

	return photo, nil
}
