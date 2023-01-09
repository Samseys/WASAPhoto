package database

func (db *appdbimpl) DeletePhoto(photoID uint64) error {
	_, err := db.c.Exec("DELETE FROM Photos WHERE id = ?", photoID)
	return err
}
