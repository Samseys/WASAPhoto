package database

func (db *appdbimpl) IdExists(id uint64) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM Users WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// return False if the requester is not banned by the other user, True otherwise
func (db *appdbimpl) IsBanned(authID uint64, otherUserID uint64) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM Bans WHERE userid = ? AND bannedid = ?)", otherUserID, authID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
