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
	err := db.c.QueryRow("SELECT Photos.id, name, ownerid, maincomment, creationdate FROM Photos INNER JOIN Users ON Photos.ownerid = Users.id WHERE Photos.id = ?", photoID).Scan(&photo.ID, &photo.Owner.Username, &photo.Owner.ID, &photo.Comment, &photo.CreationDate)
	if errors.Is(err, sql.ErrNoRows) {
		return photo, ErrPhotoNotFound
	}
	rows, err := db.c.Query("SELECT Comments.id, Users.id, Users.name, Comments.creationdate, Comments.comment FROM Comments INNER JOIN Users ON Comments.ownerid = Users.id WHERE photoid = ? ORDER BY Comments.creationdate ASC", photoID)
	if err != nil {
		return photo, err
	}

	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.ID, &comment.Owner.ID, &comment.Owner.Username, &comment.CreationDate, &comment.Comment)
		if err != nil {
			return photo, err
		}
		photo.Comments = append(photo.Comments, comment)
	}

	if err = rows.Err(); err != nil {
		return photo, err
	}

	rows.Close()

	rows, err = db.c.Query("SELECT Users.name, Likes.userid FROM Likes INNER JOIN Users on Likes.userid = Users.id WHERE photoid = ?", photoID)
	if err != nil {
		return photo, err
	}

	for rows.Next() {
		var profile UserProfileSimplified
		err = rows.Scan(&profile.Username, &profile.ID)
		if err != nil {
			return photo, err
		}
		photo.Likes = append(photo.Likes, profile)
	}

	if err = rows.Err(); err != nil {
		return photo, err
	}

	return photo, nil
}
