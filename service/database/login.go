package database

func (db *appdbimpl) Login(username string) (int, error) {
	db.c.Exec("INSERT OR IGNORE INTO Users(id, name) VALUES (NULL, ?)", username)

	var id int
	err := db.c.QueryRow("SELECT id FROM Users WHERE name=?", username).Scan(&id)

	return id, err
}
