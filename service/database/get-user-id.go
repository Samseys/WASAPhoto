package database

func (db *appdbimpl) GetUserID(username string) (uint64, error) {
	var id uint64
	err := db.c.QueryRow("SELECT id FROM Users WHERE name=?", username).Scan(&id)
	return id, err
}
