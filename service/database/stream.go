package database

func (db *appdbimpl) GetStream(userID uint64) (Stream, error) {
	var stream Stream
	rows, err := db.c.Query("SELECT id FROM (SELECT Photos.id, Photos.creationdate FROM Photos INNER JOIN Follows ON Photos.ownerid = Follows.followedid WHERE followerid = ? EXCEPT SELECT Photos.id, Photos.creationdate FROM Photos INNER JOIN Bans ON Photos.ownerid = Bans.userid WHERE bannedid = ? ORDER BY creationdate DESC)", userID, userID)
	if err != nil {
		return stream, err
	}

	for rows.Next() {
		var photoID uint64
		err = rows.Scan(&photoID)
		if err != nil {
			return stream, err
		}
		photo, err := db.GetPhotoInfoForFrontend(photoID)
		if err != nil {
			return stream, err
		}
		stream.Photos = append(stream.Photos, photo)
	}

	if err = rows.Err(); err != nil {
		return stream, err
	}

	rows.Close()

	return stream, nil
}
