package models

type User struct {
	ID			uint
	Email 		string
	Username	string
	Balance 	float32
}

func NewUser(username, email) *User {
	return &User{
		Username: username,
		Email: email
	}
}