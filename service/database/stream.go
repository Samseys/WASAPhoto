package database

func (db *appdbimpl) GetStream(userID uint64) (Stream, error) {
	var stream Stream
	rows, err := db.c.Query("SELECT Photos.id FROM Photos INNER JOIN Follows ON Photos.ownerid = Follows.followedid WHERE followerid = ? ORDER BY creationdate DESC", userID)
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

	rows.Close()

	return stream, nil
}
