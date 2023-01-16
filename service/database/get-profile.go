package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetUserProfile(userID uint64) (UserProfile, error) {
	var userProfile UserProfile
	userProfile.ID = userID
	err := db.c.QueryRow("SELECT name FROM Users WHERE id = ?", userID).Scan(&userProfile.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userProfile, ErrUserProfileNotFound
		} else {
			return userProfile, err
		}
	}

	rows, err := db.c.Query("SELECT Users.name, Follows.followerid FROM Follows INNER JOIN Users ON Users.id = Follows.followerid WHERE followedid = ?", userID)
	if err != nil {
		return userProfile, err
	}

	for rows.Next() {
		var profile UserProfileSimplified
		err = rows.Scan(&profile.Username, &profile.ID)
		if err != nil {
			return userProfile, err
		}
		userProfile.Followers = append(userProfile.Followers, profile)
	}

	if err = rows.Err(); err != nil {
		return userProfile, err
	}

	rows.Close()

	rows, err = db.c.Query("SELECT Users.name, Follows.followedid FROM Follows INNER JOIN Users ON Users.id = Follows.followedid WHERE followerid = ?", userID)
	if err != nil {
		return userProfile, err
	}
	for rows.Next() {

		var profile UserProfileSimplified
		err = rows.Scan(&profile.Username, &profile.ID)
		if err != nil {
			return userProfile, err
		}
		userProfile.Following = append(userProfile.Following, profile)
	}

	if err = rows.Err(); err != nil {
		return userProfile, err
	}

	rows.Close()

	rows, err = db.c.Query("SELECT Users.name, Bans.bannedid FROM Bans INNER JOIN Users ON Users.id = Bans.userid WHERE Bans.userid = ?", userID)
	if err != nil {
		return userProfile, err
	}

	for rows.Next() {
		var profile UserProfileSimplified
		err = rows.Scan(&profile.Username, &profile.ID)
		if err != nil {
			return userProfile, err
		}
		userProfile.Banned = append(userProfile.Banned, profile)
	}

	if err = rows.Err(); err != nil {
		return userProfile, err
	}

	rows.Close()

	rows, err = db.c.Query("SELECT Users.name, Bans.userid FROM Bans INNER JOIN Users ON Users.id = Bans.userid WHERE Bans.bannedid = ?", userID)
	if err != nil {
		return userProfile, err
	}
	for rows.Next() {

		var profile UserProfileSimplified
		err = rows.Scan(&profile.Username, &profile.ID)
		if err != nil {
			return userProfile, err
		}
		userProfile.BannedBy = append(userProfile.BannedBy, profile)
	}

	if err = rows.Err(); err != nil {
		return userProfile, err
	}

	rows.Close()

	rows, err = db.c.Query("SELECT id FROM Photos WHERE ownerid = ? ORDER BY creationdate DESC", userID)
	if err != nil {
		return userProfile, err
	}

	for rows.Next() {
		var photoID uint64
		err = rows.Scan(&photoID)
		if err != nil {
			return userProfile, err
		}
		userProfile.Photos = append(userProfile.Photos, photoID)
	}

	if err = rows.Err(); err != nil {
		return userProfile, err
	}

	rows.Close()

	return userProfile, nil
}
