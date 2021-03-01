package entity

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrImsiNotFound = errors.New("imsi not faund")
)
