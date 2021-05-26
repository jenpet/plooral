package security

type PartialPasswordSet struct {
	Password *string `json:"password"`
	PasswordConfirmation *string `json:"password_confirmation"`
}

type PasswordSet struct {
	ID int `json:"-"`
	Password string `json:"password"`
}
