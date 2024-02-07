package requests

type SignUp struct {
	FirstName            string `schema:"first_name,required"`
	LastName             string `schema:"last_name,required"`
	UserName             string `schema:"user_name,required"`
	Email                string `schema:"email,required" validate:"email"`
	SecurityQuestion     string `schema:"security_question,required"`
	Answer               string `schema:"answer,required"`
	Password             string `schema:"password,required"`
	PasswordConfirmation string `schema:"password_confirmation,required" validate:"eqfield=Password"`
}
