package database

type Name struct {
	Name string `json:"Username"`
}

type UserID struct {
	ID int `json:"UserID"`
}

type PhotoID struct {
	ID int `json:"PhotoID"`
}

type UserProfile struct {
	Username  string                  `json:"Username"`
	Photos    []string                `json:"Photos"`
	Followers []UserProfileSimplified `json:"Followers"`
	Following []UserProfileSimplified `json:"Following"`
}

type UserProfileSimplified struct {
	Username string `json:"Username"`
	ID       string `json:"UserID"`
}
