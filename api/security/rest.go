package security

type PartialUserPasswordInput struct {
	Password *string `json:"password"`
	PasswordConfirmation *string `json:"password_confirmation"`
}
