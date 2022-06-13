package service

import "errors"

var (
	DuplicateUserErr = errors.New("Duplicate user")
	UserNotFoundErr  = errors.New("User no found")
	PasswordNotMatch = errors.New("password doesn't match")
)
