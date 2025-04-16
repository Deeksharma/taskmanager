package errors

import "errors"

type Error error

var (
	ErrAuthHeaderNotPresent Error = errors.New("authorization header not present, please provide authorization header with format `Bearer {token}`")
	ErrInvalidToken         Error = errors.New("invalid token")
	ErrAuthorizationFailed  Error = errors.New("you do not have permission to perform this action")
)

var (
	ErrBadRequest Error = errors.New("bad request")
)
