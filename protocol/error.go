package protocol

type ErrorMessageType uint8

const (
	ControlMessageError ErrorMessageType = iota + 1
	DataMessageError
)

type ErrorMessage struct {
	Type               ErrorMessageType
	ControlMessageType ControlMessageType
	IndexNumber        int
	UserId             int
	RoomId             int
}

func NewErrorMessage() *ErrorMessage {
	return &ErrorMessage{}
}

func (e *ErrorMessage) Clear() {
	e.Type = 0
	e.ControlMessageType = 0
	e.IndexNumber = 0
	e.UserId = 0
	e.RoomId = 0
}

func (e *ErrorMessage) Code() []byte {
	if e.Type == ControlMessageError {
		return []byte{uint8(ErrorMessageTitle), uint8(e.Type), uint8(e.ControlMessageType)}
	}
	var (
		length    int
		index, id []byte
	)
	index = numericCoder(e.IndexNumber)
	length = 2 + len(index)
	if e.UserId != 0 {
		id = titleNumericCoder(e.UserId, uint8(userIdTitle))
		length += len(id)
	} else {
		id = titleNumericCoder(e.RoomId, uint8(roomIdTitle))
		length += len(id)
	}
	out := make([]byte, 2, length)
	copy(out, []byte{uint8(ErrorMessageTitle), uint8(e.Type)})
	out = append(out, index...)
	out = append(out, id...)
	return out
}

func (e *ErrorMessage) Decode(in []byte) error {
	if len(in) < 3 {
		return ErrorMessageTooShort
	}
	if in[0] != uint8(ErrorMessageTitle) {
		return ErrorInMessageType
	}
	var err error
	e.Type = ErrorMessageType(in[1])
	if e.Type == ControlMessageError {
		e.ControlMessageType = ControlMessageType(in[2])
		return nil
	}
	in = in[2:]
	e.IndexNumber, in, err = numericDecoder(in)
	if err != nil {
		return err
	}
	switch Title(in[0] & 0b0111_1111) {
	case userIdTitle:
		e.UserId, _, err = titleNumericDecoder(in)
		if err != nil {
			return err
		}
	case roomIdTitle:
		e.RoomId, _, err = titleNumericDecoder(in)
		if err != nil {
			return err
		}
	}
	return nil
}
