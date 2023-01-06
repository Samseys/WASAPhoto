package database

type Name struct {
	Name string `json:"Username"`
}

type UserId struct {
	Identifier int `json:"UserID"`
}

type PhotoId struct {
	Identifier int `json:"PhotoID"`
}

type UserProfile struct {
	Username  string                  `json:"Username"`
	Photos    []string                `json:"Photos"`
	Followers []UserProfileSimplified `json:"Followers"`
	Following []UserProfileSimplified `json:"Following"`
}

type UserProfileSimplified struct {
	Username   string `json:"Username"`
	Identifier string `json:"UserID"`
}
