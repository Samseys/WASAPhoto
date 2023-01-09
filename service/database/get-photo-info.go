package database

import (
	"database/sql"
	"errors"
)

var ErrPhotoNotFound = errors.New("Photo doesn't exist")

func (db *appdbimpl) GetPhotoInfo(photoID uint64) (Photo, error) {
	var photo Photo
	err := db.c.QueryRow("SELECT id, ownerid, extension, creationdate FROM Photos WHERE id = ?", photoID).Scan(&photo.ID, &photo.OwnerID, &photo.Extension, &photo.CreationDate)
	if errors.Is(err, sql.ErrNoRows) {
		return photo, ErrPhotoNotFound
	}

	return photo, nil
}
