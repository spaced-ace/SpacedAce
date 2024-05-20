package business

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Session struct {
	Id   string `json:"session"`
	User User   `json:"user"`
}
