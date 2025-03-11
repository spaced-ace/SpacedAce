package business

type User struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type Session struct {
	Id   string `json:"session"`
	User User   `json:"user"`
}
