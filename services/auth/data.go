package auth

type Account struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Salt     string `json:"-"`
	Active   bool   `json:"active"`
	UserId   string `json:"userId"`
}

type User struct {
	Id        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
