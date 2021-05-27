package security

import (
	"github.com/jenpet/plooral/errors"
	"net/http"
)

const KCredentialInputInvalid errors.Kind = "CredentialInputInvalid"

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

func (s *Service) PersistCredentials(set PartialCredentialSet) (*CredentialSet, error) {
	if *set.Password == "" || *set.Password != *set.PasswordConfirmation {
		return nil, errors.E("given passwords do not match", KCredentialInputInvalid, http.StatusBadRequest)
	}
	return s.repo.persistPassword(*set.Password)
}

func (s *Service) GenerateAndPersistPassword() (*CredentialSet, error) {
	return s.repo.persistPassword(randomSequence(generatedPasswordLength))
}

func (s *Service) VerifyCredentials(id int, password string) (bool, error) {
	return s.repo.verifyPassword(id, password)
}