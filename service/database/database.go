/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

var PhotoPath = "/home/wasa/Desktop/images/"

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Login(username string) (uint64, error)
	IdExists(id uint64) (bool, error)
	IsBanned(authID uint64, otherUserID uint64) (bool, error)
	ChangeName(username string, id uint64) error
	GetUserID(username string) (uint64, error)
	GetUserProfile(userid uint64) (UserProfile, error)
	FollowUser(followerID uint64, followedID uint64) error
	UnfollowUser(followerID uint64, followedID uint64) error
	BanUser(userID uint64, bannedID uint64) error
	UnbanUser(followerID uint64, followedID uint64) error
	UploadPhoto(userID uint64, mainComment string, extension string) (uint64, error)
	PhotoExists(photoID uint64) (bool, error)
	GetPhotoOwner(photoID uint64) (uint64, error)
	LikePhoto(userID uint64, photoID uint64) error
	UnlikePhoto(userID uint64, photoID uint64) error
	CommentPhoto(userID uint64, photoID uint64, comment string) error
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		var sqlStatements []string
		sqlStatements = append(sqlStatements, `CREATE TABLE Users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
			);`)
		sqlStatements = append(sqlStatements, `CREATE TABLE Photos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ownerid INTEGER NOT NULL,
			maincomment TEXT NOT NULL,
			extension TEXT NOT NULL,
			creationdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			FOREIGN KEY(ownerid) REFERENCES Users(id)
			);`)
		sqlStatements = append(sqlStatements, `CREATE TABLE Comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ownerid INTEGER NOT NULL,
			photoid INTEGER NOT NULL,
			creationdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			comment TEXT NOT NULL,
			FOREIGN KEY (ownerid) REFERENCES Users(id),
			FOREIGN KEY (photoid) REFERENCES Photos(id) ON DELETE CASCADE
			);`)
		sqlStatements = append(sqlStatements, `CREATE TABLE Follows (
			followerid INTEGER NOT NULL,
			followedid TEXT NOT NULL,
			FOREIGN KEY (followerid) REFERENCES Users(id),
			FOREIGN KEY (followedid) REFERENCES Users(id),
			PRIMARY KEY(followerid, followedid)
			);`)
		sqlStatements = append(sqlStatements, `CREATE TABLE Bans (
			userid INTEGER NOT NULL,
			bannedid TEXT NOT NULL,
			FOREIGN KEY (userid) REFERENCES Users(id),
			FOREIGN KEY (bannedid) REFERENCES Users(id),
			PRIMARY KEY(userid, bannedid)
			);`)
		sqlStatements = append(sqlStatements, `CREATE TABLE Likes (
			photoid INTEGER NOT NULL,
			userid TEXT NOT NULL,
			FOREIGN KEY (photoid) REFERENCES Photos(id),
			FOREIGN KEY (userid) REFERENCES Photos(id),
			PRIMARY KEY(photoid, userid)
			);`)
		for _, sqlStmt := range sqlStatements {
			_, err = db.Exec(sqlStmt)
			if err != nil {
				return nil, fmt.Errorf("error creating database structure: %w", err)
			}
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
