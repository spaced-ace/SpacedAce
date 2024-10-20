package request

type LoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type SignupForm struct {
	Email         string `form:"email"`
	Name          string `form:"name"`
	Password      string `form:"password"`
	PasswordAgain string `form:"passwordAgain"`
}
