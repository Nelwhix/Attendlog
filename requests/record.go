package requests

type Record struct {
	FirstName string `schema:"first_name,required"`
	LastName  string `schema:"last_name,required"`
	Email     string `schema:"email_address,required"`
	Signature string
}
