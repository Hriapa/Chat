package protocol

type AckMessageType uint8

const (
	Send AckMessageType = iota + 1
	Receive
	Read
)

type AckMessage struct {
	Type        AckMessageType
	IndexNumber int
	UserId      int
	Room        bool
	RoomId      int
}

func NewAckMessage() *AckMessage {
	return &AckMessage{}
}

func (a *AckMessage) Clear() {
	a.Type = 0
	a.IndexNumber = 0
	a.UserId = 0
	a.Room = false
	a.RoomId = 0
}

func (a *AckMessage) Code() []byte {
	var (
		length            int
		index, user, room []byte
	)
	index = numericCoder(a.IndexNumber)
	length = 2 + len(index)
	if a.UserId != 0 {
		user = titleNumericCoder(a.UserId, uint8(userIdTitle))
		length += len(user)
	}
	if a.Room {
		room = titleNumericCoder(a.RoomId, uint8(roomIdTitle))
		length += len(room)
	}
	out := make([]byte, 2, length)
	copy(out, []byte{uint8(AckMessageTitle), uint8(a.Type)})
	out = append(out, index...)
	if a.UserId != 0 {
		out = append(out, user...)
	}
	if a.Room {
		out = append(out, room...)
	}
	return out
}

func (a *AckMessage) Decode(in []byte) error {
	if len(in) < 5 {
		return ErrorMessageTooShort
	}
	if in[0] != byte(AckMessageTitle) {
		return ErrorInMessageType
	}
	var (
		title Title
		err   error
	)
	a.Type = AckMessageType(in[1])
	in = in[2:]
	a.IndexNumber, in, err = numericDecoder(in)
	if err != nil {
		return err
	}
	for {
		if len(in) < 2 {
			break
		}
		title = Title(in[0] & 0b0111_1111)
		switch title {
		case userIdTitle: // User Id
			a.UserId, in, err = titleNumericDecoder(in)
			if err != nil {
				return err
			}
		case roomIdTitle: // Room Id
			a.RoomId, in, err = titleNumericDecoder(in)
			if err != nil {
				return err
			}
			a.Room = true
		default:
			return nil
		}
	}
	return nil
}
