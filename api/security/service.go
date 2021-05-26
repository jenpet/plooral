package security

import (
	"github.com/jenpet/plooral/errors"
	"net/http"
)

const KPasswordsInputInvalid errors.Kind = "PasswordInputInvalid"

// generatedPasswordLength indicates how many runes should be used for system generated passwords
const generatedPasswordLength = 16

func newDefaultService() *Service {
	return &Service{repo: newDefaultRepository()}
}

func newService(r *repository) *Service {
	return &Service{repo: r}
}

type Service struct {
	repo *repository
}


func (s *Service) PersistPassword(set PartialPasswordSet) (*PasswordSet, error) {
	if *set.Password == "" || *set.Password != *set.PasswordConfirmation {
		return nil, errors.E("given passwords do not match", KPasswordsInputInvalid, http.StatusBadRequest)
	}
	return s.repo.persistPassword(*set.Password)
}

func (s *Service) GenerateAndPersistPassword() (*PasswordSet, error) {
	return s.repo.persistPassword(randomSequence(generatedPasswordLength))
}

func (s *Service) VerifyPassword(id int, password string) (bool, error) {
	return s.repo.verifyPassword(id, password)
}
