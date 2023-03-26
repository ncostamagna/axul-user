package user

import (
	"errors"
	"fmt"
	)

var NotFound = errors.New("Record not found")
var FieldIsRequired = errors.New("Required values")
var InvalidAuthentication = errors.New("Invalid authentication")
var InvalidPassword = errors.New("Invalid password")

var ErrFirstNameRequired = errors.New("first name is required")
var ErrLastNameRequired = errors.New("last name is required")
var ErrEmailRequired = errors.New("email is required")
var ErrNewPasswordRequired = errors.New("new password is required")
var ErrOldPasswordRequired = errors.New("old password is required")

type ErrNotFound struct {
	UserID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user '%s' doesn't exist", e.UserID)
}