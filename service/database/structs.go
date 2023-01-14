package database

import "time"

type Name struct {
	Name string `json:"Username"`
}

type UserID struct {
	ID uint64 `json:"UserID"`
}

type PhotoID struct {
	ID uint64 `json:"PhotoID"`
}

type Photo struct {
	ID           uint64    `json:"PhotoID"`
	OwnerID      uint64    `json:"OwnerID"`
	Extension    string    `json:"Extension"`
	CreationDate time.Time `json:"CreationDate"`
}

type Comment struct {
	ID           uint64                `json:"CommentID"`
	Comment      string                `json:"Comment"`
	Owner        UserProfileSimplified `json:"Owner"`
	CreationDate time.Time             `json:"CreationDate"`
}

type PhotoForFrontend struct {
	ID           uint64                  `json:"PhotoID"`
	Owner        UserProfileSimplified   `json:"Owner"`
	CreationDate time.Time               `json:"CreationDate"`
	Comment      string                  `json:"Comment"`
	Comments     []Comment               `json:"Comments"`
	Likes        []UserProfileSimplified `json:"Likes"`
}

type Stream struct {
	Photos []PhotoForFrontend `json:"Photos"`
}

type UserProfile struct {
	ID        uint64                  `json:"UserID"`
	Username  string                  `json:"Username"`
	Photos    []PhotoForFrontend      `json:"Photos"`
	Followers []UserProfileSimplified `json:"Followers"`
	Following []UserProfileSimplified `json:"Following"`
	Banned    []UserProfileSimplified `json:"Banned"`
	BannedBy  []UserProfileSimplified `json:"BannedBy"`
}

type UserProfileSimplified struct {
	Username string `json:"Username"`
	ID       uint64 `json:"UserID"`
}
