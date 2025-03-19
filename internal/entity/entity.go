package entity

import "errors"

type User struct {
	Username string
	Password string
}

var ErrUserExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("invalid password")
var ErrEmptyPassword = errors.New("password cannot be empty")
