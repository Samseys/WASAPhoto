package database

func (db *appdbimpl) Login(username string) (int, error) {
	db.c.Exec("INSERT OR IGNORE INTO Users(id, name) VALUES (NULL, ?)", username)

	var identifier int
	err := db.c.QueryRow("SELECT id FROM Users WHERE name=?", username).Scan(&identifier)

	return identifier, err
}
