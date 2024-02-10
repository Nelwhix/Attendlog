package requests

type Login struct {
	Email      string `schema:"email,required" validate:"email"`
	Password   string `schema:"password,required"`
	RememberMe bool   `schema:"remember_me"`
}
