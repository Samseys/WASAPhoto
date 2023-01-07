package database

func (db *appdbimpl) Login(username string) (int, error) {
	db.c.Exec("INSERT OR IGNORE INTO Users(name) VALUES (?)", username)

	var id int
	err := db.c.QueryRow("SELECT id FROM Users WHERE name=?", username).Scan(&id)

	return id, err
}
