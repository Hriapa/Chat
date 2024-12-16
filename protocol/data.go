package protocol

type DataMessageType uint8

const (
	NewMessage DataMessageType = iota + 1
	OldMessage
	UpdateMessage
)

type DataMessageFormat uint8

const (
	Text DataMessageFormat = iota + 1
	Image
	Audio
)

type Fragmentation struct {
	On           bool
	FragmentType Fragment
	Counter      uint8
}

type UserData struct {
	Fragmentation Fragmentation
	Data          []byte
}

type DataMessage struct {
	Type        DataMessageType
	Format      DataMessageFormat
	IndexNumber int
	UserId      int
	Room        bool
	RoomId      int
	UserData    UserData
}

func NewDataMessage() *DataMessage {
	return &DataMessage{}
}

func (d *DataMessage) Clear() {
	d.Type = 0
	d.Format = 0
	d.IndexNumber = 0
	d.UserId = 0
	d.Room = false
	d.RoomId = 0
	d.UserData.clear()
}

func (d *DataMessage) Code() []byte {
	var (
		length                  int
		index, user, room, data []byte
	)
	index = numericCoder(d.IndexNumber)
	length = 3 + len(index)
	if d.UserId != 0 {
		user = titleNumericCoder(d.UserId, uint8(userIdTitle))
		length += len(user)
	}
	if d.Room {
		room = titleNumericCoder(d.RoomId, uint8(roomIdTitle))
		length += len(room)
	}
	data = d.UserData.dataCoder()
	length += len(data)
	out := make([]byte, 3, length)
	copy(out, []byte{uint8(DataMessageTitle), uint8(d.Type), uint8(d.Format)})
	out = append(out, index...)
	if d.UserId != 0 {
		out = append(out, user...)
	}
	if d.Room {
		out = append(out, room...)
	}
	out = append(out, data...)
	return out
}

func (d *DataMessage) Decode(in []byte) error {
	if len(in) < 5 {
		return ErrorMessageTooShort
	}
	if in[0] != byte(DataMessageTitle) {
		return ErrorInMessageType
	}
	var (
		title Title
		err   error
	)
	d.Type = DataMessageType(in[1])
	d.Format = DataMessageFormat(in[2])
	in = in[3:]
	d.IndexNumber, in, err = numericDecoder(in)
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
			d.UserId, in, err = titleNumericDecoder(in)
			if err != nil {
				return err
			}
		case roomIdTitle: // Room Id
			d.RoomId, in, err = titleNumericDecoder(in)
			if err != nil {
				return err
			}
			d.Room = true
		case userDataTitle: //User Data
			in, err = d.UserData.dataDecoder(in)
			if err != nil {
				return err
			}
		default:
			return nil
		}
	}
	return nil
}

func (d *UserData) clear() {
	d.Fragmentation.On = false
	d.Fragmentation.FragmentType = 4
	d.Fragmentation.Counter = 0
	d.Data = nil
}

func (u *UserData) dataCoder() []byte {
	var (
		fragment uint8
		length   int
		data     []byte
	)
	length = len(u.Data) + 1
	if u.Fragmentation.On {
		length += 1
		fragment = uint8(u.Fragmentation.FragmentType) | u.Fragmentation.Counter
	}
	data = dataCoder(u.Data)
	length += len(data)
	out := make([]byte, 1, length)
	copy(out, []byte{uint8(userDataTitle)})
	if u.Fragmentation.On {
		out[0] = out[0] | 0b1000_0000 // set fr bit 1
		out = append(out, fragment)
	}
	out = append(out, data...)
	return out
}

func (u *UserData) dataDecoder(in []byte) ([]byte, error) {
	var (
		err error
	)
	if len(in) < 2 {
		return in, ErrorMessageTooShort
	}
	u.Fragmentation.On = (in[0] >> 7) == 1
	in = in[1:]
	if u.Fragmentation.On {
		u.Fragmentation.FragmentType = Fragment(in[0] & 0b1100_0000)
		u.Fragmentation.Counter = in[0] & 0b0011_1111
		in = in[1:]
	}
	if len(in) == 0 {
		return in, ErrorDecode
	}
	u.Data, in, err = dataDecoder(in)
	return in, err
}
