package rest

import (
	"github.com/jenpet/plooral/errors"
	"net/http"
)

const (
	KUserInputInvalid errors.Kind = "UserInputInvalid"
)

func WrapUserInputInvalidError(err error) error {
	return errors.Ef("invalid user input", http.StatusBadRequest, KUserInputInvalid, err)
}