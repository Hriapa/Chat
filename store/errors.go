package store

import "errors"

var (
	ErrorConnectToDatabase = errors.New("error: no connect to database")
	ErrorLoginNotFound     = errors.New("error: login not found")
	ErrorPasswordIncorrect = errors.New("error: password incorrect")
	ErrorUserNotFound      = errors.New("error: user not found")
	ErrorNoMessage         = errors.New("error: no message in database")
)
