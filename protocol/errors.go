package protocol

import "errors"

var (
	ErrorInMessageType   = errors.New(`error, incorrect message type`)
	ErrorLengthDecode    = errors.New(`error, length in decode message incorrect`)
	ErrorTitleIncorrect  = errors.New(`error, incorrect title value`)
	ErrorMessageTooShort = errors.New(`error, message too short for decoding`)
	ErrorDecode          = errors.New(`error, cann't decode message`)
)
