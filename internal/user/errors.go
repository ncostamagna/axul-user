package user

import "errors"

var NotFound = errors.New("Record not found")
var FieldIsRequired = errors.New("Required values")
var InvalidAuthentication = errors.New("Invalid authentication")
