package database

func (db *appdbimpl) BanUser(userID uint64, bannedID uint64) error {
	sqlres, err := db.c.Exec("INSERT INTO Bans (userid, bannedid) VALUES (?, ?) ON CONFLICT DO NOTHING", userID, bannedID)
	if err != nil {
		return err
	}

	rowsAffected, _ := sqlres.RowsAffected()
	if rowsAffected == 0 {
		return ErrAlreadyBanned
	}

	return nil
}
