package database

func (db *appdbimpl) PhotoExists(photoID uint64) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM Photos WHERE id = ?)", photoID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
