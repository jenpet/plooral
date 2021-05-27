package security

type PartialCredentialSet struct {
	Password *string `json:"password"`
	PasswordConfirmation *string `json:"password_confirmation"`
}

type CredentialSet struct {
	ID int `json:"-"`
	Password string `json:"password"`
}

func (cs *CredentialSet) IsValid() bool {
	return cs.HasValidID() && len(cs.Password) > 0
}

func (cs *CredentialSet) HasValidID() bool {
	return cs != nil && cs.ID > 0
}
