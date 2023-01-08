package database

func (db *appdbimpl) Login(username string) (uint64, error) {
	sqlRes, err := db.c.Exec("INSERT INTO Users(name) VALUES (?) ON CONFLICT DO NOTHING", username)
	if err != nil {
		return 0, err
	}
	affectedRows, err := sqlRes.RowsAffected()
	if err != nil {
		return 0, err
	}

	var id uint64
	if affectedRows == 0 {
		err = db.c.QueryRow("SELECT id FROM Users WHERE name=?", username).Scan(&id)
		return id, err
	} else {
		lastinsertid, err := sqlRes.LastInsertId()
		return uint64(lastinsertid), err
	}
}
