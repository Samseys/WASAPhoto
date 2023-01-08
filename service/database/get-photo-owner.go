package database

func (db *appdbimpl) GetPhotoOwner(photoID uint64) (uint64, error) {
	var ownerID uint64
	err := db.c.QueryRow("SELECT ownerid FROM Photos WHERE id = ?", photoID).Scan(&ownerID)
	if err != nil {
		return 0, err
	}
	return ownerID, nil
}
