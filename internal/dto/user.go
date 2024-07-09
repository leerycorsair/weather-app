package dto

type DTOSignUp struct {
	Login    string `json:"login"  db:"login"`
	Password string `json:"password"  db:"password"`
	Email    string `json:"email"  db:"email"`
}

type DTOSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
