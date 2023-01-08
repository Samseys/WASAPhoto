package database

type Name struct {
	Name string `json:"Username"`
}

type UserID struct {
	ID uint64 `json:"UserID"`
}

type PhotoID struct {
	ID uint64 `json:"PhotoID"`
}

type Comment struct {
	Comment string `json:"Comment"`
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
