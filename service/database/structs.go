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
	ID           uint64                `json:"PhotoID"`
	Owner        UserProfileSimplified `json:"Owner"`
	CreationDate time.Time             `json:"CreationDate"`
	Comments     []Comment             `json:"Comments"`
}

type Stream struct {
	Photos []PhotoForFrontend `json:"Photos"`
}

type UserProfile struct {
	Username  string                  `json:"Username"`
	Photos    []uint64                `json:"Photos"`
	Followers []UserProfileSimplified `json:"Followers"`
	Following []UserProfileSimplified `json:"Following"`
}

type UserProfileSimplified struct {
	Username string `json:"Username"`
	ID       uint64 `json:"UserID"`
}
