package models

// User represents the user model
// @Description User model
type User struct {
	Id       int    `json:"id"  db:"id"`             // @Description User ID
	Login    string `json:"login"  db:"login"`       // @Description User login
	Password string `json:"password"  db:"password"` // @Description User password
	Email    string `json:"email"  db:"email"`       // @Description User email
}
