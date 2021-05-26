package security

type PartialCredentialSet struct {
	Password *string `json:"password"`
	PasswordConfirmation *string `json:"password_confirmation"`
}

type CredentialSet struct {
	ID int `json:"-"`
	Password string `json:"password"`
}
