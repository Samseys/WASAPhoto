package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetUserProfile(userid int) (UserProfile, error) {
	var userProfile UserProfile
	err := db.c.QueryRow("SELECT name FROM Users WHERE id = ?", userid).Scan(&userProfile.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return userProfile, err
	}

	rows, err := db.c.Query("SELECT Users.name, Follows.followerid FROM Follows INNER JOIN Users ON Users.id = Follows.followerid WHERE followedid = ?", userid)
	if err != nil {
		return userProfile, err
	}

	followers := []UserProfileSimplified{}
	for rows.Next() {
		var profile UserProfileSimplified
		rows.Scan(&profile.Username, &profile.Identifier)
		if err != nil {
			return userProfile, err
		}
		followers = append(followers, profile)
	}

	userProfile.Followers = followers

	rows, err = db.c.Query("SELECT Users.name, Follows.followedid FROM Follows INNER JOIN Users ON Users.id = Follows.followedid WHERE followerid = ?", userid)
	if err != nil {
		return userProfile, err
	}

	following := []UserProfileSimplified{}
	for rows.Next() {
		var profile UserProfileSimplified
		err = rows.Scan(&profile.Username, &profile.Identifier)
		if err != nil {
			return userProfile, err
		}
		following = append(following, profile)

	}

	userProfile.Following = following
	fmt.Printf("userProfile.Following: %v\n", following)

	rows, err = db.c.Query("SELECT id FROM Photos WHERE ownerid = ? ORDER BY timedate DESC", userid)
	if err != nil {
		return userProfile, err
	}

	photos := []string{}
	for rows.Next() {
		var photoid string
		rows.Scan(&photoid)
		if err != nil {
			return userProfile, err
		}
		photos = append(photos, photoid)
	}

	userProfile.Photos = photos
	return userProfile, nil
}
