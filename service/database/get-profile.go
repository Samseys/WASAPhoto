package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetUserProfile(userID uint64) (UserProfile, error) {
	var userProfile UserProfile
	err := db.c.QueryRow("SELECT name FROM Users WHERE id = ?", userID).Scan(&userProfile.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return userProfile, ErrUserProfileNotFound
	}

	rows, err := db.c.Query("SELECT Users.name, Follows.followerid FROM Follows INNER JOIN Users ON Users.id = Follows.followerid WHERE followedid = ?", userID)
	if err != nil {
		return userProfile, err
	}

	followers := []UserProfileSimplified{}
	for rows.Next() {
		var profile UserProfileSimplified
		err = rows.Scan(&profile.Username, &profile.ID)
		if err != nil {
			return userProfile, err
		}
		followers = append(followers, profile)
	}

	rows.Close()

	userProfile.Followers = followers

	rows, err = db.c.Query("SELECT Users.name, Follows.followedid FROM Follows INNER JOIN Users ON Users.id = Follows.followedid WHERE followerid = ?", userID)
	if err != nil {
		return userProfile, err
	}

	following := []UserProfileSimplified{}
	for rows.Next() {
		var profile UserProfileSimplified
		err = rows.Scan(&profile.Username, &profile.ID)
		if err != nil {
			return userProfile, err
		}
		following = append(following, profile)
	}

	rows.Close()

	userProfile.Following = following

	rows, err = db.c.Query("SELECT id FROM Photos WHERE ownerid = ? ORDER BY creationdate DESC", userID)
	if err != nil {
		return userProfile, err
	}

	photos := []uint64{}
	for rows.Next() {
		var photoid uint64
		err = rows.Scan(&photoid)
		if err != nil {
			return userProfile, err
		}
		photos = append(photos, photoid)
	}

	rows.Close()

	userProfile.Photos = photos
	return userProfile, nil
}
