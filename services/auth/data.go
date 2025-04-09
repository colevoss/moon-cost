package auth

type Account struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type User struct {
	Id        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
