package role

import (
	"errors"
	"fmt"
)

var ErrUserIDAndAppAreRequired = errors.New("user id and app are required")

/*var FieldIsRequired = errors.New("Required values")
var InvalidAuthentication = errors.New("Invalid authentication")
var InvalidPassword = errors.New("Invalid password")

var ErrFirstNameRequired = errors.New("first name is required")
var ErrLastNameRequired = errors.New("last name is required")
var ErrEmailRequired = errors.New("email is required")
var ErrNewPasswordRequired = errors.New("new password is required")
var ErrOldPasswordRequired = errors.New("old password is required")*/

type ErrUserAppNotFound struct {
	UserID string
	App    string
}

func (e ErrUserAppNotFound) Error() string {
	return fmt.Sprintf("user '%s' with '%s' app doesn't exist", e.UserID, e.App)
}

type InvalidRole struct {
	Role string
}

func (e InvalidRole) Error() string {
	return fmt.Sprintf("the '%s' isn't valid", e.Role)
}
