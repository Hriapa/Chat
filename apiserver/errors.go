package apiserver

import "errors"

var (
	errIncorrectLoginOrPassword = errors.New("login or password incorrect")
	errConnectToDatabase        = errors.New("error connect to database")
	errNotAuthenticated         = errors.New("not authenticated")
	errUserNotExist             = errors.New("user not exist")
	errUserNotFound             = errors.New("user not found")
	errIncorrectPassword        = errors.New("password incorrect")
)
